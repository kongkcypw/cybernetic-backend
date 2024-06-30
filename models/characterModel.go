package models

// Character represents the character model
type Character struct {
	Id            string `bson:"_id" json:"id"`
	UserId        string `json:"userId" validate:"required"`
	CharacterName string `json:"characterName" validate:"required"`
}
