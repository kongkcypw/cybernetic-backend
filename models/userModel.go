package models

// Character represents the character model
type UserCharacter struct {
	Id            string `bson:"_id" json:"id"`
	UserId        string `bson:"userId" json:"userId" validate:"required"`
	CharacterName string `bson:"characterName" json:"characterName"`
	HeighestLevel int    `bson:"heighestLevel" json:"heighestLevel"`
}

type UserGamePlayLevel struct {
	UserId      string `json:"userId"`
	LevelPlayed []int  `json:"levelPlayed"`
}
