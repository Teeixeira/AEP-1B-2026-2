package models

import "go.mongodb.org/mongo-driver/v2/bson"

const TipoPonto = "Point"

type Localizacao struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type Crime struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Tipo        string        `bson:"tipo" json:"tipo"`
	Descricao   string        `bson:"descricao" json:"descricao"`
	DataHora    string        `bson:"data_hora" json:"data_hora"`
	Localizacao Localizacao   `bson:"localizacao" json:"localizacao"`
}
