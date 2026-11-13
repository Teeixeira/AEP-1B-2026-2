package migrations

import (
	"context"
	"fmt"
	"log"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	migrate.MustRegister(func(ctx context.Context, db *mongo.Database) error {
		collection := db.Collection("crimes")

		// 1. Encontrar documentos no formato antigo com localizacao.coordinates ou com campo localizacao_geo residual
		filter := bson.M{
			"$or": []bson.M{
				{"localizacao.coordinates": bson.M{"$exists": true}},
				{"localizacao.type": bson.M{"$exists": true}},
				{"localizacao_geo": bson.M{"$exists": true}},
				{"localizacao.latitude": bson.M{"$exists": false}},
			},
		}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("erro ao buscar crimes para migração: %w", err)
		}
		defer cursor.Close(ctx)

		migratedCount := 0
		for cursor.Next(ctx) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				return fmt.Errorf("erro ao decodificar documento durante migração: %w", err)
			}

			id := doc["_id"]
			locRaw := doc["localizacao"]
			locGeoRaw := doc["localizacao_geo"]

			var lng, lat float64
			var coords []float64

			if locRaw != nil {
				coords = extractCoords(locRaw)
			}
			if len(coords) < 2 && locGeoRaw != nil {
				coords = extractCoords(locGeoRaw)
			}

			if len(coords) >= 2 {
				lng = coords[0]
				lat = coords[1]
			} else if locRaw != nil {
				var ok bool
				lat, lng, ok = extractLatLng(locRaw)
				if !ok {
					continue
				}
			} else {
				continue
			}

			// Define longitude como primeiro campo (X) e latitude como segundo (Y) no BSON.
			// Sobrescrever "localizacao" já substitui o objeto por inteiro, eliminando coordinates e type sem conflito.
			update := bson.M{
				"$set": bson.M{
					"localizacao": bson.D{
						{Key: "longitude", Value: lng},
						{Key: "latitude", Value: lat},
					},
				},
				"$unset": bson.M{
					"localizacao_geo": "",
				},
			}

			if _, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update); err != nil {
				return fmt.Errorf("erro ao atualizar documento %v na migração: %w", id, err)
			}
			migratedCount++
		}

		if err := cursor.Err(); err != nil {
			return fmt.Errorf("erro durante iteração de documentos: %w", err)
		}

		log.Printf("[MIGRATION 1] %d registros convertidos diretamente para localizacao {longitude, latitude}.", migratedCount)

		// 2. Remove eventual índice em localizacao_geo e garante índice 2dsphere em localizacao
		dropOldGeoIndexIfExists(ctx, collection, "localizacao_geo")
		dropOldGeoIndexIfExists(ctx, collection, "localizacao")

		indexModel := mongo.IndexModel{
			Keys: bson.D{{Key: "localizacao", Value: "2dsphere"}},
		}
		if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
			log.Printf("[MIGRATION 1] Aviso: falha ao criar índice localizacao_2dsphere: %v", err)
		}

		return nil
	}, func(ctx context.Context, db *mongo.Database) error {
		collection := db.Collection("crimes")

		filter := bson.M{
			"localizacao.latitude": bson.M{"$exists": true},
		}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("erro ao buscar crimes para rollback: %w", err)
		}
		defer cursor.Close(ctx)

		revertedCount := 0
		for cursor.Next(ctx) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				return fmt.Errorf("erro ao decodificar documento no rollback: %w", err)
			}

			id := doc["_id"]
			locRaw, ok := doc["localizacao"]
			if !ok || locRaw == nil {
				continue
			}

			lat, lng, okCoord := extractLatLng(locRaw)
			if !okCoord {
				continue
			}

			update := bson.M{
				"$set": bson.M{
					"localizacao": bson.M{
						"type":        "Point",
						"coordinates": []float64{lng, lat},
					},
				},
				"$unset": bson.M{
					"localizacao_geo": "",
				},
			}

			if _, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update); err != nil {
				return fmt.Errorf("erro ao reverter documento %v: %w", id, err)
			}
			revertedCount++
		}

		log.Printf("[MIGRATION 1] %d registros revertidos para formato GeoJSON.", revertedCount)

		dropOldGeoIndexIfExists(ctx, collection, "localizacao")

		indexModel := mongo.IndexModel{
			Keys: bson.D{{Key: "localizacao", Value: "2dsphere"}},
		}
		if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
			log.Printf("[MIGRATION 1] Aviso: falha ao recriar índice localizacao_2dsphere: %v", err)
		}

		return nil
	})
}

func extractCoords(locRaw interface{}) []float64 {
	var coordsRaw interface{}
	switch m := locRaw.(type) {
	case bson.M:
		coordsRaw = m["coordinates"]
	case map[string]interface{}:
		coordsRaw = m["coordinates"]
	case bson.D:
		for _, elem := range m {
			if elem.Key == "coordinates" {
				coordsRaw = elem.Value
				break
			}
		}
	}
	if coordsRaw == nil {
		return nil
	}

	var coords []float64
	switch v := coordsRaw.(type) {
	case bson.A:
		for _, item := range v {
			if num, ok := toFloat(item); ok {
				coords = append(coords, num)
			}
		}
	case []interface{}:
		for _, item := range v {
			if num, ok := toFloat(item); ok {
				coords = append(coords, num)
			}
		}
	case []float64:
		coords = v
	}
	return coords
}

func extractLatLng(locRaw interface{}) (float64, float64, bool) {
	var latRaw, lngRaw interface{}
	switch m := locRaw.(type) {
	case bson.M:
		latRaw = m["latitude"]
		lngRaw = m["longitude"]
	case map[string]interface{}:
		latRaw = m["latitude"]
		lngRaw = m["longitude"]
	case bson.D:
		for _, elem := range m {
			if elem.Key == "latitude" {
				latRaw = elem.Value
			}
			if elem.Key == "longitude" {
				lngRaw = elem.Value
			}
		}
	}
	latVal, okLat := toFloat(latRaw)
	lngVal, okLng := toFloat(lngRaw)
	if !okLat || !okLng {
		return 0, 0, false
	}
	return latVal, lngVal, true
}

func dropOldGeoIndexIfExists(ctx context.Context, collection *mongo.Collection, keyName string) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var idx bson.M
		if err := cursor.Decode(&idx); err != nil {
			continue
		}
		name, _ := idx["name"].(string)
		keyRaw, ok := idx["key"]
		if !ok {
			continue
		}

		isTarget := false
		if keyD, isD := keyRaw.(bson.D); isD {
			for _, elem := range keyD {
				if elem.Key == keyName && elem.Value == "2dsphere" {
					isTarget = true
					break
				}
			}
		} else if keyMap, isM := keyRaw.(bson.M); isM {
			if keyMap[keyName] == "2dsphere" {
				isTarget = true
			}
		}

		if isTarget && name != "" {
			_ = collection.Indexes().DropOne(ctx, name)
		}
	}
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}
