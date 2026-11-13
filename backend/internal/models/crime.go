package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Localizacao armazena latitude e longitude em graus decimais.
type Localizacao struct {
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
}

type Crime struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Tipo        string        `bson:"tipo" json:"tipo"`
	Descricao   string        `bson:"descricao" json:"descricao"`
	DataHora    time.Time     `bson:"data_hora" json:"data_hora"`
	Localizacao Localizacao   `bson:"localizacao" json:"localizacao"`
}
