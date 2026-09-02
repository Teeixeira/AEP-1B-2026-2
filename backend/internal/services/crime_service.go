package services

import (
	"context"
	"errors"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrRaioInvalido = errors.New("raio de busca deve ser maior que zero")

type CrimeServiceInterface interface {
	Create(ctx context.Context, crime *models.Crime) error
	FindAll(ctx context.Context) ([]models.Crime, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*models.Crime, error)
	FindNear(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error)
	Update(ctx context.Context, id bson.ObjectID, crime models.Crime) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

type CrimeService struct {
	repository repositories.CrimeRepositoryInterface
}

func NewCrimeService(repository repositories.CrimeRepositoryInterface) *CrimeService {
	return &CrimeService{repository: repository}
}

func (service *CrimeService) Create(
	ctx context.Context,
	crime *models.Crime,
) error {
	crime.Localizacao.Type = models.TipoPonto
	return service.repository.Create(ctx, crime)
}

func (service *CrimeService) FindAll(ctx context.Context) ([]models.Crime, error) {
	return service.repository.FindAll(ctx)
}

func (service *CrimeService) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Crime, error) {
	return service.repository.FindByID(ctx, id)
}

func (service *CrimeService) FindNear(
	ctx context.Context,
	longitude float64,
	latitude float64,
	raioMetros float64,
) ([]models.Crime, error) {
	if raioMetros <= 0 {
		return nil, ErrRaioInvalido
	}
	return service.repository.FindNear(ctx, longitude, latitude, raioMetros)
}

func (service *CrimeService) Update(
	ctx context.Context,
	id bson.ObjectID,
	crime models.Crime,
) error {
	crime.Localizacao.Type = models.TipoPonto
	return service.repository.Update(ctx, id, crime)
}

func (service *CrimeService) Delete(ctx context.Context, id bson.ObjectID) error {
	return service.repository.Delete(ctx, id)
}
