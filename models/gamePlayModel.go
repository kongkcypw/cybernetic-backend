package models

type LevelDescription struct {
	Head    string `bson:"head" json:"head"`
	Content string `bson:"content" json:"content"`
}

type GamePlayLevel struct {
	LevelId     string             `bson:"_id" json:"id"`
	LevelNumber int64              `bson:"levelNumber" json:"levelNumber"`
	LevelName   string             `bson:"levelName" json:"levelName"`
	Description []LevelDescription `bson:"description" json:"description"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
}
