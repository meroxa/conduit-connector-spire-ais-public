package ais

func vesselQuery() string {
	return `
	query ($first: Int!, $after: String){
	        vessels(first:$first, after:$after) {
				pageInfo {
				 hasNextPage
				 endCursor
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
