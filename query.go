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

func vesselQuery() string {
	return `
	query ($first: Int!, $after: String){
	        vessels(first:$first, after:$after, lastPositionUpdate: { startTime: "2023-11-12T21:00:48.768Z" }) {
				pageInfo {
				 hasNextPage
				 endCursor
			   }
			   totalCount { #recordset details
				value
				relation
			   }
			   nodes {
				 id
				 updateTimestamp
				 staticData {
				   aisClass
				   flag
				   name
				   callsign
				   timestamp
				   updateTimestamp
				   shipType
				   shipSubType
				   mmsi
				   imo
				   callsign
				   dimensions {
					 a
					 b
					 c
					 d
					 width
					 length
				   }
				 }
				 lastPositionUpdate {
				   accuracy
				   collectionType
				   course
				   heading
				   latitude
				   longitude
				   maneuver
				   navigationalStatus
				   rot
				   speed
				   timestamp
				   updateTimestamp
				 }
				 currentVoyage {
				   destination
				   draught
				   eta
				   timestamp
				   updateTimestamp
				 }
			   }
			 }
	    }
	`
}
