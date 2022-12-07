package model

type LoyaltyConfig struct {
	Formula    []int32 `json:"formula,omitempty" bson:"formula,omitempty"`
	PointCycle string  `json:"point_cycle,omitempty" bson:"point_cycle,omitempty"`
	Ranks      []Rank  `json:"ranks,omitempty" bson:"ranks,omitempty"`
}

type Rank struct {
	Title        string `json:"title,omitempty" bson:"title,omitempty"`
	MinimumScore int32  `json:"minimum_score" bson:"minimum_score"`
}
