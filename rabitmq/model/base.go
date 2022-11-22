package model

import (
	"time"
)

// Mgostream full document stream
type Mgostream struct {
	OPType string `bson:"operationType" json:"event_type"`
	Key    struct {
		ID string `bson:"_id" json:"id"`
	} `bson:"documentKey" json:"id"`
	Detail struct {
		UpdateFiled map[string]interface{} `bson:"updatedFields" json:"update_filed"`
		RemoveFiled []string               `bson:"removedFields" json:"remove_filed"`
	} `bson:"updateDescription" json:"detail"`
	FullDoc map[string]interface{} `bson:"fullDocument" json:"full_doc"`
	NS      struct {
		DB   string `bson:"db" json:"db"`
		Coll string `bson:"coll" json:"coll"`
	} `bson:"ns" json:"ns"`
	EventTime time.Time `bson:"clusterTime" json:"event_time"`
}