// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
