package repositories

import (
	"context"
	"fmt"

	"AEP-1B-2026-2/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CrimeRepository struct {
	collection *mongo.Collection
}

func NewCrimeRepository(database *mongo.Database) *CrimeRepository {
	return &CrimeRepository{
		collection: database.Collection("crimes"),
	}
}

func (repository *CrimeRepository) EnsureGeoIndex(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "localizacao", Value: "2dsphere"}},
	}
	_, err := repository.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("erro ao criar índice geoespacial: %w", err)
	}
	return nil
}

func (repository *CrimeRepository) Create(
	ctx context.Context,
	crime *models.Crime,
) error {
	crime.ID = bson.NewObjectID()
	_, err := repository.collection.InsertOne(ctx, crime)
	if err != nil {
		return fmt.Errorf("erro ao inserir crime: %w", err)
	}
	return nil
}

func (repository *CrimeRepository) FindAll(
	ctx context.Context,
) ([]models.Crime, error) {
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar crimes: %w", err)
	}
	defer cursor.Close(ctx)

	var crimes []models.Crime
	if err := cursor.All(ctx, &crimes); err != nil {
		return nil, fmt.Errorf("erro ao decodificar crimes: %w", err)
	}
	return crimes, nil
}

func (repository *CrimeRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Crime, error) {
	var crime models.Crime
	err := repository.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&crime)
	if err != nil {
		return nil, err
	}
	return &crime, nil
}

func (repository *CrimeRepository) FindNear(
	ctx context.Context,
	longitude float64,
	latitude float64,
	raioMetros float64,
) ([]models.Crime, error) {
	filter := bson.M{
		"localizacao": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        models.TipoPonto,
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": raioMetros,
			},
		},
	}

	cursor, err := repository.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar crimes próximos: %w", err)
	}
	defer cursor.Close(ctx)

	var crimes []models.Crime
	if err := cursor.All(ctx, &crimes); err != nil {
		return nil, fmt.Errorf("erro ao decodificar crimes próximos: %w", err)
	}
	return crimes, nil
}

func (repository *CrimeRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	crime models.Crime,
) error {
	update := bson.M{
		"$set": bson.M{
			"tipo":        crime.Tipo,
			"descricao":   crime.Descricao,
			"data_hora":   crime.DataHora,
			"localizacao": crime.Localizacao,
		},
	}
	result, err := repository.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar crime: %w", err)
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repository *CrimeRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	result, err := repository.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("erro ao excluir crime: %w", err)
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
