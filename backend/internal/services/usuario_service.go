package services

import (
	"context"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UsuarioService struct {
	repository *repositories.UsuarioRepository
}

func NewUsuarioService(
	repository *repositories.UsuarioRepository,
) *UsuarioService {

	return &UsuarioService{
		repository: repository,
	}
}

func (service *UsuarioService) Create(
	ctx context.Context,
	usuario *models.Usuario,
) error {

	return service.repository.Create(ctx, usuario)
}

func (service *UsuarioService) FindAll(
	ctx context.Context,
) ([]models.Usuario, error) {

	return service.repository.FindAll(ctx)
}

func (service *UsuarioService) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Usuario, error) {

	return service.repository.FindByID(ctx, id)
}

func (service *UsuarioService) Update(
	ctx context.Context,
	id bson.ObjectID,
	usuario models.Usuario,
) error {

	return service.repository.Update(ctx, id, usuario)
}

func (service *UsuarioService) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {

	return service.repository.Delete(ctx, id)
}
