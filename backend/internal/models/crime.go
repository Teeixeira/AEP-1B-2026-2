package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Localizacao armazena coordenadas geográficas no formato de pares legado (longitude, latitude)
type Localizacao struct {
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
}

type Crime struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Tipo        string        `bson:"tipo" json:"tipo"`
	Descricao   string        `bson:"descricao" json:"descricao"`
	DataHora    string        `bson:"data_hora" json:"data_hora"`
	Localizacao Localizacao   `bson:"localizacao" json:"localizacao"`
}
