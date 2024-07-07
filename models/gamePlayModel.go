package models

type GamePlayLevel struct {
	LevelId     string `bson:"_id" json:"id"`
	LevelNumber int64  `bson:"levelNumber" json:"levelNumber"`
	LevelName   string `bson:"levelName" json:"levelName"`
	Description string `bson:"description" json:"description"`
	IsActive    bool   `bson:"isActive" json:"isActive"`
}
