package domain

type RoadType string

const (
	RoadCycleTrack   RoadType = "cycle_track"
	RoadBikeLane     RoadType = "bike_lane"
	RoadResidential  RoadType = "residential"
	RoadCommonStreet RoadType = "common_street"
	RoadBusyAvenue   RoadType = "busy_avenue"
	RoadHighway      RoadType = "highway"
	RoadUnsafe       RoadType = "unsafe"
)

type RouteProfile string

const (
	ProfileBalanced RouteProfile = "balanced"
	ProfileSafest   RouteProfile = "safest"
	ProfileFastest  RouteProfile = "fastest"
)

type Node struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Lat  float64 `json:"lat,omitempty"`
	Lon  float64 `json:"lon,omitempty"`
	Zone string  `json:"zone"`
}

type Edge struct {
	ID            string   `json:"id"`
	From          string   `json:"from"`
	To            string   `json:"to"`
	Name          string   `json:"name"`
	Distance      float64  `json:"distance"`
	StreetType    RoadType `json:"street_type"`
	SafetyScore   float64  `json:"safety_score"`
	SpeedLimit    int      `json:"speed_limit"`
	TrafficLevel  int      `json:"traffic_level"`
	Bidirectional bool     `json:"bidirectional"`
}

type CityMap struct {
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Route struct {
	OriginID      string  `json:"origin_id"`
	DestinationID string  `json:"destination_id"`
	Profile       string  `json:"profile"`
	Nodes         []Node  `json:"nodes"`
	Edges         []Edge  `json:"edges"`
	TotalDistance float64 `json:"total_distance"`
	EstimatedTime float64 `json:"estimated_time"`
	SafetyScore   float64 `json:"safety_score"`
}

type CyclistPosition struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	NodeHint string  `json:"node_hint,omitempty"`
}
