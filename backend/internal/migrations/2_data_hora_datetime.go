package migrations

import (
	"context"
	"fmt"
	"time"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	migrate.MustRegister(migrateDataHoraToDate, rollbackDataHoraToString)
}

func migrateDataHoraToDate(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("crimes")
	cursor, err := collection.Find(ctx, bson.M{
		"data_hora": bson.M{"$type": "string"},
	})
	if err != nil {
		return fmt.Errorf("erro ao buscar crimes para conversão de data_hora: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc struct {
			ID       bson.ObjectID `bson:"_id"`
			DataHora string        `bson:"data_hora"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("erro ao decodificar crime durante conversão de data_hora: %w", err)
		}

		date, err := parseDataHora(doc.DataHora)
		if err != nil {
			return fmt.Errorf("data_hora inválida no crime %v (%q): %w", doc.ID, doc.DataHora, err)
		}

		if _, err := collection.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{
			"$set": bson.M{"data_hora": date},
		}); err != nil {
			return fmt.Errorf("erro ao converter data_hora do crime %v: %w", doc.ID, err)
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("erro durante iteração de crimes para conversão de data_hora: %w", err)
	}
	return nil
}

func parseDataHora(value string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339Nano, value)
	if err == nil {
		return date, nil
	}

	// Legacy timestamps without a timezone are interpreted as UTC.
	return time.ParseInLocation("2006-01-02T15:04:05.999999999", value, time.UTC)
}

func rollbackDataHoraToString(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("crimes")
	cursor, err := collection.Find(ctx, bson.M{
		"data_hora": bson.M{"$type": "date"},
	})
	if err != nil {
		return fmt.Errorf("erro ao buscar crimes para rollback de data_hora: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc struct {
			ID       bson.ObjectID `bson:"_id"`
			DataHora time.Time     `bson:"data_hora"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("erro ao decodificar crime durante rollback de data_hora: %w", err)
		}

		value := doc.DataHora.UTC().Format(time.RFC3339Nano)

		if _, err := collection.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{
			"$set": bson.M{"data_hora": value},
		}); err != nil {
			return fmt.Errorf("erro ao reverter data_hora do crime %v: %w", doc.ID, err)
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("erro durante iteração de crimes para rollback de data_hora: %w", err)
	}
	return nil
}
