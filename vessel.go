package ais

type Vessels struct {
	PageInfo PageInfo `json:"pageInfo"`
	Nodes    []Node   `json:"nodes"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type Node struct {
	ID                 string             `json:"id"`
	UpdateTimestamp    string             `json:"updateTimestamp"`
	StaticData         StaticData         `json:"staticData"`
	LastPositionUpdate LastPositionUpdate `json:"lastPositionUpdate"`
	CurrentVoyage      CurrentVoyage      `json:"currentVoyage"`
}

// This is unused as we are not returning structured data. We are instead returning []byte.
//func (n Node) toStructuredData() sdk.StructuredData {
//	return sdk.StructuredData{
//		"id":                 n.ID,
//		"updateTimestamp":    n.UpdateTimestamp,
//		"staticData":         n.StaticData,
//		"lastPositionUpdate": n.LastPositionUpdate,
//		"currentVoyage":      n.CurrentVoyage,
//	}
//}

type StaticData struct {
	AisClass        string     `json:"aisClass"`
	Flag            string     `json:"flag"`
	Name            string     `json:"name"`
	Callsign        string     `json:"callsign"`
	Timestamp       string     `json:"timestamp"`
	UpdateTimestamp string     `json:"updateTimestamp"`
	ShipType        string     `json:"shipType"`
	ShipSubType     string     `json:"shipSubType"`
	MMSI            int        `json:"mmsi"`
	IMO             int        `json:"imo"`
	Dimensions      Dimensions `json:"dimensions"`
}

type Dimensions struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	C      float64 `json:"c"`
	D      float64 `json:"d"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
}

type LastPositionUpdate struct {
	Accuracy           string  `json:"accuracy"`
	CollectionType     string  `json:"collectionType"`
	Course             float64 `json:"course"`
	Heading            float64 `json:"heading"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Maneuver           string  `json:"maneuver"`
	NavigationalStatus string  `json:"navigationalStatus"`
	Rot                float64 `json:"rot"`
	Speed              float64 `json:"speed"`
	Timestamp          string  `json:"timestamp"`
	UpdateTimestamp    string  `json:"updateTimestamp"`
}

type CurrentVoyage struct {
	Destination     string  `json:"destination"`
	Draught         float64 `json:"draught"`
	ETA             string  `json:"eta"`
	Timestamp       string  `json:"timestamp"`
	UpdateTimestamp string  `json:"updateTimestamp"`
}
