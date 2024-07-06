package models

// Character represents the character model
type UserCharacter struct {
	Id            string `bson:"_id" json:"id"`
	UserId        string `json:"userId" validate:"required"`
	CharacterName string `json:"characterName" validate:"required"`
}

type UserGamePlayLevel struct {
	UserId      string `json:"userId"`
	LevelPlayed []int  `json:"levelPlayed"`
}
