#!/usr/bin/env python3
import argparse
import json
import math
import xml.etree.ElementTree as ET
from collections import defaultdict

ALLOWED_HIGHWAYS = {
    "motorway", "trunk", "primary", "secondary", "tertiary", "unclassified",
    "residential", "living_street", "service", "cycleway"
}
SKIP_SERVICE = {"driveway", "drive-through", "parking_aisle", "emergency_access"}
SKIP_ACCESS = {"private", "no"}

TYPE_DEFAULTS = {
    "cycle_track": (0.94, 25, 1),
    "bike_lane": (0.86, 30, 2),
    "residential": (0.76, 30, 1),
    "common_street": (0.66, 35, 2),
    "busy_avenue": (0.46, 50, 4),
    "highway": (0.16, 80, 5),
    "unsafe": (0.22, 60, 5),
}

HIGHWAY_LABELS = {
    "motorway": "Rodovia sem nome",
    "trunk": "Rodovia sem nome",
    "primary": "Avenida principal sem nome",
    "secondary": "Avenida sem nome",
    "tertiary": "Via coletora sem nome",
    "unclassified": "Rua sem nome",
    "residential": "Rua residencial sem nome",
    "living_street": "Rua compartilhada sem nome",
    "service": "Acesso local sem nome",
    "cycleway": "Ciclovia sem nome",
}

def parse_maxspeed(value, fallback):
    if not value:
        return fallback
    digits = "".join(ch for ch in value if ch.isdigit())
    if not digits:
        return fallback
    speed = int(digits)
    return max(10, min(speed, 110))

def has_cycle_value(tags, values):
    for key, value in tags.items():
        if key == "cycleway" or key.startswith("cycleway:"):
            if any(v in value for v in values):
                return True
    return False

def road_type(tags):
    highway = tags.get("highway", "")
    if highway == "cycleway" or tags.get("bicycle") == "designated" or has_cycle_value(tags, {"track", "separate"}):
        return "cycle_track"
    if has_cycle_value(tags, {"lane", "shared_lane", "opposite_lane"}):
        return "bike_lane"
    if highway in {"motorway", "trunk"}:
        return "highway"
    if highway in {"primary", "secondary", "tertiary"}:
        return "busy_avenue"
    if highway in {"residential", "living_street"}:
        return "residential"
    return "common_street"

def defaults_for(street_type, tags):
    safety, speed, traffic = TYPE_DEFAULTS[street_type]
    speed = parse_maxspeed(tags.get("maxspeed"), speed)
    if tags.get("surface") in {"unpaved", "gravel", "dirt", "ground", "sand"}:
        safety = max(0.1, safety - 0.12)
    if tags.get("bicycle") == "no":
        return "unsafe", 0.16, speed, 5
    return street_type, round(safety, 2), speed, traffic

def haversine(a, b):
    radius = 6371000.0
    lat1, lon1 = math.radians(a[0]), math.radians(a[1])
    lat2, lon2 = math.radians(b[0]), math.radians(b[1])
    dlat = lat2 - lat1
    dlon = lon2 - lon1
    h = math.sin(dlat / 2) ** 2 + math.cos(lat1) * math.cos(lat2) * math.sin(dlon / 2) ** 2
    return 2 * radius * math.asin(math.sqrt(h))

def load_osm(paths):
    nodes = {}
    ways = {}
    for path in paths:
        root = ET.parse(path).getroot()
        for node in root.findall("node"):
            nodes[node.attrib["id"]] = (float(node.attrib["lat"]), float(node.attrib["lon"]))
        for way in root.findall("way"):
            tags = {tag.attrib["k"]: tag.attrib["v"] for tag in way.findall("tag")}
            highway = tags.get("highway")
            if highway not in ALLOWED_HIGHWAYS:
                continue
            if tags.get("access") in SKIP_ACCESS or tags.get("motor_vehicle") == "private":
                continue
            if highway == "service" and tags.get("service") in SKIP_SERVICE and not tags.get("name"):
                continue
            refs = [nd.attrib["ref"] for nd in way.findall("nd") if nd.attrib["ref"] in nodes]
            if len(refs) < 2:
                continue
            ways[way.attrib["id"]] = {"id": way.attrib["id"], "tags": tags, "refs": refs}
    return nodes, ways

def normalize(nodes):
    lats = [lat for lat, _ in nodes.values()]
    lons = [lon for _, lon in nodes.values()]
    min_lat, max_lat = min(lats), max(lats)
    min_lon, max_lon = min(lons), max(lons)
    def scale(ref):
        lat, lon = nodes[ref]
        x = 4 + ((lon - min_lon) / (max_lon - min_lon)) * 96
        y = 4 + ((max_lat - lat) / (max_lat - min_lat)) * 88
        return round(x, 4), round(y, 4)
    return scale

def bidirectional(tags):
    if tags.get("oneway:bicycle") == "no" or "opposite" in tags.get("cycleway", ""):
        return True
    return tags.get("oneway") not in {"yes", "true", "1"}

def node_zone(types):
    if "highway" in types:
        return "highway"
    if "busy_avenue" in types:
        return "arterial"
    if "cycle_track" in types or "bike_lane" in types or "common_street" in types:
        return "collector"
    return "residential"

def largest_component(used, edges):
    adjacency = defaultdict(set)
    id_to_ref = {f"N{ref}": ref for ref in used}
    for edge in edges:
        a = id_to_ref[edge["from"]]
        b = id_to_ref[edge["to"]]
        adjacency[a].add(b)
        adjacency[b].add(a)

    seen = set()
    components = []
    for ref in used:
        if ref in seen:
            continue
        stack = [ref]
        seen.add(ref)
        component = set()
        while stack:
            current = stack.pop()
            component.add(current)
            for nxt in adjacency[current]:
                if nxt not in seen:
                    seen.add(nxt)
                    stack.append(nxt)
        components.append(component)

    if not components:
        return used, edges
    keep = max(components, key=len)
    filtered_edges = [edge for edge in edges if edge["from"][1:] in keep and edge["to"][1:] in keep]
    return keep, filtered_edges


def build(paths):
    raw_nodes, ways = load_osm(paths)
    used = set()
    node_names = defaultdict(set)
    node_types = defaultdict(set)
    edge_defs = []
    edge_seen = set()

    for way in sorted(ways.values(), key=lambda w: int(w["id"])):
        tags = way["tags"]
        name = tags.get("name") or tags.get("ref") or HIGHWAY_LABELS.get(tags.get("highway"), "Via sem nome")
        base_type = road_type(tags)
        street_type, safety, speed, traffic = defaults_for(base_type, tags)
        refs = way["refs"]
        for idx, (a, b) in enumerate(zip(refs, refs[1:])):
            if a == b or a not in raw_nodes or b not in raw_nodes:
                continue
            distance = haversine(raw_nodes[a], raw_nodes[b])
            if distance < 1:
                continue
            key = (a, b, name, street_type)
            if key in edge_seen:
                continue
            edge_seen.add(key)
            used.add(a); used.add(b)
            node_names[a].add(name); node_names[b].add(name)
            node_types[a].add(street_type); node_types[b].add(street_type)
            edge_defs.append({
                "id": f"OSM{way['id']}_{idx}",
                "from": f"N{a}",
                "to": f"N{b}",
                "name": name,
                "distance": round(distance, 2),
                "street_type": street_type,
                "safety_score": safety,
                "speed_limit": speed,
                "traffic_level": traffic,
                "bidirectional": bidirectional(tags),
            })

    used, edge_defs = largest_component(used, edge_defs)

    scale = normalize({ref: raw_nodes[ref] for ref in used})
    nodes = []
    for ref in sorted(used, key=int):
        names = sorted(n for n in node_names[ref] if n)
        if len(names) >= 2:
            name = " / ".join(names[:2])
        elif len(names) == 1:
            name = names[0]
        else:
            name = f"Ponto OSM {ref}"
        x, y = scale(ref)
        lat, lon = raw_nodes[ref]
        nodes.append({"id": f"N{ref}", "name": name, "x": x, "y": y, "lat": round(lat, 7), "lon": round(lon, 7), "zone": node_zone(node_types[ref])})

    return {"name": "Pelotas RS - grafo OSM Rei-Julian", "nodes": nodes, "edges": edge_defs}

def main():
    parser = argparse.ArgumentParser(description="Converte recortes OSM XML para o mapa do Rei-Julian.")
    parser.add_argument("inputs", nargs="+", help="Arquivos .osm XML baixados do OpenStreetMap")
    parser.add_argument("-o", "--output", default="maps/city_01.json")
    args = parser.parse_args()
    city = build(args.inputs)
    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(city, f, ensure_ascii=False, indent=2)
        f.write("\n")
    print(f"wrote {args.output}: {len(city['nodes'])} nodes, {len(city['edges'])} edges")

if __name__ == "__main__":
    main()
