const canvas = document.querySelector("#map");
const ctx = canvas.getContext("2d");
const origin = document.querySelector("#origin");
const destination = document.querySelector("#destination");
const profile = document.querySelector("#profile");
const summary = document.querySelector("#summary");
const directions = document.querySelector("#directions");
const tracking = document.querySelector("#tracking");
const statsBox = document.querySelector("#stats");
const mapName = document.querySelector("#map-name");
const hint = document.querySelector("#hint");
const pickOrigin = document.querySelector("#pick-origin");
const pickDestination = document.querySelector("#pick-destination");
const zoomIn = document.querySelector("#zoom-in");
const zoomOut = document.querySelector("#zoom-out");
const resetViewButton = document.querySelector("#reset-view");
const zoomLabel = document.querySelector("#zoom-label");
const originSearch = document.querySelector("#origin-search");
const destinationSearch = document.querySelector("#destination-search");
const originResults = document.querySelector("#origin-results");
const destinationResults = document.querySelector("#destination-results");
const comparison = document.querySelector("#comparison");
const compareButton = document.querySelector("#compare");
const gpsButton = document.querySelector("#gps");
const stopSimulationButton = document.querySelector("#stop-simulation");

const colors = {
  cycle_track: "#16824b",
  bike_lane: "#1f6f8b",
  residential: "#8794a6",
  common_street: "#5f6b7a",
  busy_avenue: "#c96f1f",
  highway: "#161f2e",
  unsafe: "#b42318"
};

const streetLabels = {
  cycle_track: "ciclovia",
  bike_lane: "ciclofaixa",
  residential: "rua residencial",
  common_street: "rua comum",
  busy_avenue: "avenida movimentada",
  highway: "rodovia",
  unsafe: "trecho inseguro"
};

let city;
let stats;
let route;
let selectedMode = "origin";
let hoveredNode;
let nodesByID = {};
let degrees = new Map();
let selectableNodes = [];
let viewport = { zoom: 1, offsetX: 0, offsetY: 0 };
let isDragging = false;
let dragStart = { x: 0, y: 0, offsetX: 0, offsetY: 0 };
let pointerDownPoint;
let sessionID = "";
let simulationTimer;
let simulatedPosition;
const searchCache = { origin: new Map(), destination: new Map() };

function worldPoint(node) {
  return { x: 70 + node.x * 10.5, y: 60 + node.y * 7.5 };
}

function screenPointFromWorld(point) {
  return {
    x: point.x * viewport.zoom + viewport.offsetX,
    y: point.y * viewport.zoom + viewport.offsetY
  };
}

function worldPointFromScreen(point) {
  return {
    x: (point.x - viewport.offsetX) / viewport.zoom,
    y: (point.y - viewport.offsetY) / viewport.zoom
  };
}

function scale(node) {
  return screenPointFromWorld(worldPoint(node));
}

function updateZoomLabel() {
  zoomLabel.textContent = `${Math.round(viewport.zoom * 100)}%`;
}

function setZoom(nextZoom, anchor = { x: canvas.width / 2, y: canvas.height / 2 }) {
  const oldZoom = viewport.zoom;
  const zoom = Math.max(0.55, Math.min(8, nextZoom));
  if (zoom === oldZoom) return;
  const world = worldPointFromScreen(anchor);
  viewport.zoom = zoom;
  viewport.offsetX = anchor.x - world.x * zoom;
  viewport.offsetY = anchor.y - world.y * zoom;
  updateZoomLabel();
  draw();
}

function resetView() {
  viewport = { zoom: 1, offsetX: 0, offsetY: 0 };
  updateZoomLabel();
  draw();
}

function normalizeText(value) {
  return value.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
}

function formatKm(meters) {
  return `${(meters / 1000).toFixed(1)} km`;
}

function formatSafety(score) {
  return `${Math.round(score * 100)}%`;
}

function renderStats() {
  statsBox.innerHTML = `
    <div class="metric"><strong>${stats.node_count.toLocaleString("pt-BR")}</strong><span>pontos</span></div>
    <div class="metric"><strong>${stats.edge_count.toLocaleString("pt-BR")}</strong><span>trechos</span></div>
    <div class="metric"><strong>${formatKm(stats.cycling_distance)}</strong><span>malha cicloviaria</span></div>
    <div class="metric"><strong>${formatSafety(stats.average_safety)}</strong><span>seguranca media</span></div>
  `;
}

function buildIndexes() {
  nodesByID = Object.fromEntries(city.nodes.map((node) => [node.id, node]));
  degrees = new Map(city.nodes.map((node) => [node.id, 0]));
  for (const edge of city.edges) {
    degrees.set(edge.from, (degrees.get(edge.from) || 0) + 1);
    degrees.set(edge.to, (degrees.get(edge.to) || 0) + 1);
  }
  selectableNodes = city.nodes
    .filter((node) => isSelectable(node))
    .sort((a, b) => scoreNode(b) - scoreNode(a) || a.name.localeCompare(b.name, "pt-BR"))
    .slice(0, 1800)
    .sort((a, b) => a.name.localeCompare(b.name, "pt-BR"));
}

function isGenericAccess(node) {
  const name = normalizeText(node.name);
  return name.includes("acesso local sem nome") || name === "rua sem nome" || name.startsWith("ponto osm");
}

function isSelectable(node) {
  const name = normalizeText(node.name);
  if (isGenericAccess(node)) return false;
  if (name.includes("sem nome") && (degrees.get(node.id) || 0) < 5) return false;
  return (degrees.get(node.id) || 0) >= 3 || node.zone !== "residential";
}

function scoreNode(node) {
  const degree = degrees.get(node.id) || 0;
  const named = normalizeText(node.name).includes("sem nome") ? 0 : 10;
  const zone = node.zone === "arterial" ? 6 : node.zone === "collector" ? 4 : node.zone === "highway" ? 2 : 0;
  return degree * 3 + named + zone;
}

function optionLabel(node) {
  return `${node.name} (${node.id})`;
}

function ensureOption(select, node) {
  if (!node || Array.from(select.options).some((option) => option.value === node.id)) return;
  select.add(new Option(optionLabel(node), node.id), 0);
}

function fillSelectors() {
  origin.textContent = "";
  destination.textContent = "";
  for (const node of selectableNodes) {
    origin.add(new Option(optionLabel(node), node.id));
    destination.add(new Option(optionLabel(node), node.id));
  }
  const defaultOrigin = findNode(["andrade neves", "coronel gomes carneiro"]) || selectableNodes[0];
  const defaultDestination = findNode(["rio grande do sul", "antonio augusto de assumpcao"]) || selectableNodes[1];
  ensureOption(origin, defaultOrigin);
  ensureOption(destination, defaultDestination);
  origin.value = defaultOrigin.id;
  destination.value = defaultDestination.id;
}

function edgeVisible(edge, routeEdgeIDs) {
  if (routeEdgeIDs.has(edge.id)) return true;
  if (viewport.zoom < 0.85) {
    return edge.street_type === "highway" || edge.street_type === "busy_avenue" || edge.street_type === "cycle_track" || edge.street_type === "bike_lane";
  }
  if (viewport.zoom < 1.6) {
    return edge.street_type !== "residential" || edge.distance > 90;
  }
  if (viewport.zoom < 2.4) {
    return edge.street_type !== "common_street" || edge.distance > 25;
  }
  return true;
}

function findNode(terms) {
  const normalizedTerms = terms.map(normalizeText);
  return city.nodes.find((node) => {
    const name = normalizeText(node.name);
    return normalizedTerms.every((term) => name.includes(term));
  });
}

function draw() {
  if (!city) return;
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  ctx.lineCap = "round";
  ctx.lineJoin = "round";
  const routeEdgeIDs = new Set((route?.edges || []).map((edge) => edge.id));
  ctx.save();

  for (const edge of city.edges) {
    if (!edgeVisible(edge, routeEdgeIDs)) continue;
    const a = scale(nodesByID[edge.from]);
    const b = scale(nodesByID[edge.to]);
    ctx.strokeStyle = colors[edge.street_type] || "#555";
    ctx.lineWidth = Math.max(1, (edge.street_type === "highway" ? 7 : edge.street_type === "cycle_track" ? 4 : edge.street_type === "bike_lane" ? 3.5 : 2) * Math.min(1.4, Math.sqrt(viewport.zoom)));
    ctx.globalAlpha = routeEdgeIDs.has(edge.id) ? 0.18 : edge.street_type === "highway" ? 0.26 : 0.52;
    ctx.beginPath();
    ctx.moveTo(a.x, a.y);
    ctx.lineTo(b.x, b.y);
    ctx.stroke();
  }

  ctx.globalAlpha = 1;
  if (route) {
    for (const edge of route.edges) {
      const a = scale(nodesByID[edge.from]);
      const b = scale(nodesByID[edge.to]);
      ctx.strokeStyle = "#6f42c1";
      ctx.lineWidth = Math.max(5, 7 * Math.min(1.5, Math.sqrt(viewport.zoom)));
      ctx.beginPath();
      ctx.moveTo(a.x, a.y);
      ctx.lineTo(b.x, b.y);
      ctx.stroke();
    }
  }

  drawMarker(nodesByID[origin.value], "#16824b", 9);
  drawMarker(nodesByID[destination.value], "#b42318", 9);
  if (hoveredNode && hoveredNode.id !== origin.value && hoveredNode.id !== destination.value && !isGenericAccess(hoveredNode)) {
    drawMarker(hoveredNode, "#6f42c1", 7);
  }
  if (simulatedPosition) {
    drawPosition(simulatedPosition, "#0f766e", 8);
  }
  drawLabels();
  ctx.restore();
}

function drawPosition(pos, color, radius) {
  const p = screenPointFromWorld({ x: 70 + pos.x * 10.5, y: 60 + pos.y * 7.5 });
  ctx.fillStyle = color;
  ctx.strokeStyle = "#ffffff";
  ctx.lineWidth = 3;
  ctx.beginPath();
  ctx.arc(p.x, p.y, radius, 0, Math.PI * 2);
  ctx.fill();
  ctx.stroke();
}

function drawMarker(node, color, radius) {
  if (!node) return;
  const p = scale(node);
  ctx.fillStyle = "#ffffff";
  ctx.strokeStyle = color;
  ctx.lineWidth = Math.max(3, 4 * Math.min(1.3, Math.sqrt(viewport.zoom)));
  ctx.beginPath();
  ctx.arc(p.x, p.y, radius, 0, Math.PI * 2);
  ctx.fill();
  ctx.stroke();
}

function drawLabels() {
  const important = new Set([origin.value, destination.value]);
  if (hoveredNode) important.add(hoveredNode.id);
  ctx.font = "12px system-ui, sans-serif";
  ctx.textBaseline = "middle";
  for (const id of important) {
    const node = nodesByID[id];
    if (!node) continue;
    const p = scale(node);
    const text = node.name.length > 70 ? `${node.name.slice(0, 67)}...` : node.name;
    const width = Math.min(ctx.measureText(text).width + 12, 430);
    ctx.fillStyle = "rgba(255,255,255,0.94)";
    ctx.strokeStyle = "rgba(38,52,69,0.18)";
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.roundRect(p.x + 11, p.y - 11, width, 22, 5);
    ctx.fill();
    ctx.stroke();
    ctx.fillStyle = "#17202a";
    ctx.fillText(text, p.x + 17, p.y, width - 10);
  }
}

function canvasPoint(event) {
  const rect = canvas.getBoundingClientRect();
  return {
    x: ((event.clientX - rect.left) / rect.width) * canvas.width,
    y: ((event.clientY - rect.top) / rect.height) * canvas.height
  };
}

function nearestNode(point) {
  let best;
  let bestDistance = Infinity;
  for (const node of city.nodes) {
    if (isGenericAccess(node) && viewport.zoom < 4) continue;
    const p = worldPoint(node);
    const distance = Math.hypot(p.x - point.x, p.y - point.y);
    if (distance < bestDistance) {
      best = node;
      bestDistance = distance;
    }
  }
  if (best && isGenericAccess(best) && viewport.zoom < 4) return undefined;
  return bestDistance <= 20 / viewport.zoom ? best : undefined;
}

function setMode(mode) {
  selectedMode = mode;
  pickOrigin.setAttribute("aria-pressed", String(mode === "origin"));
  pickDestination.setAttribute("aria-pressed", String(mode === "destination"));
  hint.textContent = mode === "origin" ? "Clique em um ponto para alterar a origem." : "Clique em um ponto para alterar o destino.";
}

async function loadMap() {
  const [mapPayload, statsPayload] = await Promise.all([
    fetch("/map").then((r) => r.json()),
    fetch("/stats").then((r) => r.json())
  ]);
  city = mapPayload;
  stats = statsPayload;
  mapName.textContent = city.name;
  buildIndexes();
  renderStats();
  fillSelectors();
  updateZoomLabel();
  draw();
}

async function calculateRoute() {
  if (origin.value === destination.value) {
    summary.textContent = "Origem e destino precisam ser diferentes.";
    return;
  }
  const response = await fetch("/routes", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: sessionID, origin_id: origin.value, destination_id: destination.value, profile: profile.value })
  });
  const payload = await response.json();
  if (!response.ok) {
    summary.textContent = payload.error || "Erro ao calcular rota.";
    return;
  }
  sessionID = payload.session_id || sessionID;
  route = payload.route;
  simulatedPosition = undefined;
  comparison.textContent = "";
  renderRoute();
  tracking.textContent = "Tracking pronto. Use Simular desvio para testar recálculo automático.";
  draw();
}

function routeSteps(edges) {
  const steps = [];
  for (const edge of edges) {
    const last = steps[steps.length - 1];
    if (last && last.name === edge.name && last.street_type === edge.street_type) {
      last.distance += edge.distance;
    } else {
      steps.push({ name: edge.name, street_type: edge.street_type, distance: edge.distance });
    }
  }
  return steps;
}

function renderRoute() {
  const start = route.nodes[0];
  const end = route.nodes[route.nodes.length - 1];
  summary.innerHTML = `
    <strong>${start.name}</strong> ate <strong>${end.name}</strong><br>
    Distancia: ${formatKm(route.total_distance)}<br>
    Tempo estimado: ${route.estimatedTime?.toFixed?.(1) || route.estimated_time.toFixed(1)} min<br>
    Seguranca media: ${formatSafety(route.safety_score)}
  `;
  directions.innerHTML = `
    <strong>Instrucoes</strong>
    <ol>${routeSteps(route.edges).slice(0, 80).map((edge) => `<li>Siga por ${edge.name} (${streetLabels[edge.street_type] || edge.street_type}) por ${formatKm(edge.distance)}.</li>`).join("")}</ol>
  `;
}

async function compareProfiles() {
  if (origin.value === destination.value) {
    comparison.textContent = "Origem e destino precisam ser diferentes.";
    return;
  }
  const response = await fetch("/routes/compare", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ origin_id: origin.value, destination_id: destination.value })
  });
  const payload = await response.json();
  if (!response.ok) {
    comparison.textContent = payload.error || "Erro ao comparar rotas.";
    return;
  }
  const rows = ["fastest", "balanced", "safest"].map((key) => {
    const item = payload.routes[key];
    return `<li><strong>${key}</strong>: ${formatKm(item.total_distance)}, ${item.estimated_time.toFixed(1)} min, ${formatSafety(item.safety_score)} seguranca</li>`;
  }).join("");
  comparison.innerHTML = `<strong>Comparacao</strong><ul>${rows}</ul>`;
}

async function startSimulation() {
  if (!route || !sessionID) {
    tracking.textContent = "Calcule uma rota antes de simular movimento.";
    return;
  }
  stopSimulation();
  const start = await fetch("/simulation/start", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: sessionID, step_meters: 140 })
  });
  if (!start.ok) {
    const payload = await start.json();
    tracking.textContent = payload.error || "Erro ao iniciar simulacao.";
    return;
  }
  tracking.textContent = "Simulacao em andamento.";
  simulationTimer = window.setInterval(stepSimulation, 550);
}

async function stepSimulation() {
  const response = await fetch("/simulation/step", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: sessionID })
  });
  const payload = await response.json();
  if (!response.ok) {
    tracking.textContent = payload.error || "Erro na simulacao.";
    stopSimulation();
    return;
  }
  simulatedPosition = payload.position;
  tracking.innerHTML = `Simulando perto de <strong>${payload.node.name}</strong>.`;
  if (payload.state.done) {
    tracking.textContent = "Simulacao concluida.";
    stopSimulation(false);
  }
  draw();
}

function stopSimulation(clearPosition = false) {
  if (simulationTimer) window.clearInterval(simulationTimer);
  simulationTimer = undefined;
  if (clearPosition) {
    simulatedPosition = undefined;
    draw();
  }
}

async function useGPS() {
  if (!navigator.geolocation) {
    tracking.textContent = "GPS indisponivel neste navegador.";
    return;
  }
  navigator.geolocation.getCurrentPosition(async (position) => {
    const node = nearestGeoNode(position.coords.latitude, position.coords.longitude);
    if (!node) {
      tracking.textContent = "Nao foi possivel relacionar o GPS ao mapa.";
      return;
    }
    ensureOption(origin, node);
    origin.value = node.id;
    tracking.innerHTML = `GPS associado a <strong>${node.name}</strong>.`;
    draw();
    if (route && sessionID) {
      const response = await fetch("/tracking/update", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ session_id: sessionID, x: node.x, y: node.y, node_hint: node.id })
      });
      const payload = await response.json();
      if (payload.new_route) {
        route = payload.new_route;
        renderRoute();
        draw();
      }
    }
  }, () => {
    tracking.textContent = "Permissao de GPS negada ou indisponivel.";
  });
}

function nearestGeoNode(lat, lon) {
  let best;
  let bestDistance = Infinity;
  for (const node of city.nodes) {
    if (!node.lat || !node.lon || isGenericAccess(node)) continue;
    const distance = Math.hypot((node.lat - lat) * 111000, (node.lon - lon) * 95500);
    if (distance < bestDistance) {
      best = node;
      bestDistance = distance;
    }
  }
  return bestDistance < 5000 ? best : undefined;
}

async function searchNode(kind, query) {
  const datalist = kind === "origin" ? originResults : destinationResults;
  const cache = searchCache[kind];
  cache.clear();
  datalist.textContent = "";
  if (query.trim().length < 3) return;
  const response = await fetch(`/search?q=${encodeURIComponent(query)}&limit=12`);
  const payload = await response.json();
  for (const node of payload.results || []) {
    const label = optionLabel(node);
    cache.set(label, node);
    datalist.append(new Option(label, label));
  }
}

function applySearch(kind, value) {
  const cache = searchCache[kind];
  const node = cache.get(value);
  if (!node) return;
  const select = kind === "origin" ? origin : destination;
  ensureOption(select, node);
  select.value = node.id;
  route = undefined;
  draw();
}

canvas.addEventListener("wheel", (event) => {
  event.preventDefault();
  const factor = event.deltaY < 0 ? 1.18 : 1 / 1.18;
  setZoom(viewport.zoom * factor, canvasPoint(event));
}, { passive: false });

canvas.addEventListener("pointerdown", (event) => {
  const point = canvasPoint(event);
  isDragging = true;
  pointerDownPoint = point;
  dragStart = { x: point.x, y: point.y, offsetX: viewport.offsetX, offsetY: viewport.offsetY };
  canvas.setPointerCapture(event.pointerId);
});

canvas.addEventListener("pointermove", (event) => {
  const point = canvasPoint(event);
  if (isDragging) {
    viewport.offsetX = dragStart.offsetX + point.x - dragStart.x;
    viewport.offsetY = dragStart.offsetY + point.y - dragStart.y;
    canvas.style.cursor = "grabbing";
    draw();
    return;
  }
  const next = nearestNode(worldPointFromScreen(point));
  if (next?.id === hoveredNode?.id) return;
  hoveredNode = next;
  canvas.style.cursor = hoveredNode ? "pointer" : "grab";
  draw();
});

canvas.addEventListener("pointerup", (event) => {
  const point = canvasPoint(event);
  const moved = pointerDownPoint && Math.hypot(point.x - pointerDownPoint.x, point.y - pointerDownPoint.y) > 4;
  isDragging = false;
  canvas.releasePointerCapture(event.pointerId);
  canvas.style.cursor = hoveredNode ? "pointer" : "grab";
  if (moved) return;
  const node = nearestNode(worldPointFromScreen(point));
  if (!node) return;
  if (selectedMode === "origin") {
    ensureOption(origin, node);
    origin.value = node.id;
  } else {
    ensureOption(destination, node);
    destination.value = node.id;
  }
  route = undefined;
  directions.textContent = "";
  summary.textContent = "Ponto atualizado. Calcule a rota.";
  draw();
});

canvas.addEventListener("pointercancel", () => {
  isDragging = false;
  canvas.style.cursor = hoveredNode ? "pointer" : "grab";
});

canvas.addEventListener("mouseleave", () => {
  hoveredNode = undefined;
  isDragging = false;
  draw();
});

zoomIn.addEventListener("click", () => setZoom(viewport.zoom * 1.25));
zoomOut.addEventListener("click", () => setZoom(viewport.zoom / 1.25));
resetViewButton.addEventListener("click", resetView);

origin.addEventListener("change", draw);
destination.addEventListener("change", draw);
document.querySelector("#route").addEventListener("click", calculateRoute);
document.querySelector("#simulate").addEventListener("click", startSimulation);
stopSimulationButton.addEventListener("click", () => stopSimulation(true));
compareButton.addEventListener("click", compareProfiles);
gpsButton.addEventListener("click", useGPS);
originSearch.addEventListener("input", () => searchNode("origin", originSearch.value));
destinationSearch.addEventListener("input", () => searchNode("destination", destinationSearch.value));
originSearch.addEventListener("change", () => applySearch("origin", originSearch.value));
destinationSearch.addEventListener("change", () => applySearch("destination", destinationSearch.value));
pickOrigin.addEventListener("click", () => setMode("origin"));
pickDestination.addEventListener("click", () => setMode("destination"));

loadMap().catch((error) => {
  summary.textContent = `Erro ao carregar o mapa: ${error.message}`;
});
