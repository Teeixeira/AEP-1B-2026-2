package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Usuario struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Nome  string        `bson:"nome" json:"nome"`
	Email string        `bson:"email" json:"email"`
	Senha string        `bson:"senha" json:"senha"`
	Ativo bool          `bson:"ativo" json:"ativo"`
}
