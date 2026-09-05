package repositories

import (
	"context"

	"AEP-1B-2026-2/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UsuarioRepository struct {
	collection *mongo.Collection
}

type UsuarioRepositoryInterface interface {
	Create(ctx context.Context, usuario *models.Usuario) error
	FindAll(ctx context.Context) ([]models.Usuario, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*models.Usuario, error)
	Update(ctx context.Context, id bson.ObjectID, usuario models.Usuario) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

func NewUsuarioRepository(database *mongo.Database) *UsuarioRepository {
	return &UsuarioRepository{
		collection: database.Collection("usuarios"),
	}
}

func (repository *UsuarioRepository) Create(
	ctx context.Context,
	usuario *models.Usuario,
) error {

	usuario.ID = bson.NewObjectID()

	_, err := repository.collection.InsertOne(ctx, usuario)

	if err != nil {
		return err
	}

	return nil
}

func (repository *UsuarioRepository) FindAll(
	ctx context.Context,
) ([]models.Usuario, error) {

	cursor, err := repository.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var usuarios []models.Usuario

	err = cursor.All(ctx, &usuarios)

	if err != nil {
		return nil, err
	}

	return usuarios, nil
}

func (repository *UsuarioRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Usuario, error) {

	var usuario models.Usuario

	err := repository.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&usuario)

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

func (repository *UsuarioRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	usuario models.Usuario,
) error {

	update := bson.M{
		"$set": bson.M{
			"nome":  usuario.Nome,
			"email": usuario.Email,
			"senha": usuario.Senha,
			"ativo": usuario.Ativo,
		},
	}

	result, err := repository.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (repository *UsuarioRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	result, err := repository.collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
