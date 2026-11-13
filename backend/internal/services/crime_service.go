package services

import (
	"context"
	"errors"
	"time"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	ErrRaioInvalido       = errors.New("raio de busca deve ser maior que zero")
	ErrCamposObrigatorios = errors.New("campos obrigatórios não preenchidos: tipo, descricao e data_hora são necessários")
	ErrLatitudeInvalida   = errors.New("latitude deve estar entre -90 e 90")
	ErrLongitudeInvalida  = errors.New("longitude deve estar entre -180 e 180")
	ErrPeriodoInvalido    = errors.New("data_inicio deve ser anterior ou igual a data_fim")
)

type CrimeServiceInterface interface {
	Create(ctx context.Context, crime *models.Crime) error
	FindAll(ctx context.Context) ([]models.Crime, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*models.Crime, error)
	FindNear(ctx context.Context, longitude, latitude, raioMetros float64, tipo string, dataInicio, dataFim time.Time) ([]models.Crime, error)
	Update(ctx context.Context, id bson.ObjectID, crime models.Crime) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

type CrimeService struct {
	repository repositories.CrimeRepositoryInterface
}

func NewCrimeService(repository repositories.CrimeRepositoryInterface) *CrimeService {
	return &CrimeService{repository: repository}
}

func validarCrime(crime *models.Crime) error {
	if crime.Tipo == "" || crime.Descricao == "" || crime.DataHora.IsZero() {
		return ErrCamposObrigatorios
	}

	lat := crime.Localizacao.Latitude
	lng := crime.Localizacao.Longitude

	if lat < -90 || lat > 90 {
		return ErrLatitudeInvalida
	}
	if lng < -180 || lng > 180 {
		return ErrLongitudeInvalida
	}

	return nil
}

func (service *CrimeService) Create(
	ctx context.Context,
	crime *models.Crime,
) error {
	if err := validarCrime(crime); err != nil {
		return err
	}
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
	tipo string,
	dataInicio time.Time,
	dataFim time.Time,
) ([]models.Crime, error) {
	if raioMetros <= 0 {
		return nil, ErrRaioInvalido
	}
	if !dataInicio.IsZero() && !dataFim.IsZero() && dataInicio.After(dataFim) {
		return nil, ErrPeriodoInvalido
	}
	return service.repository.FindNear(ctx, longitude, latitude, raioMetros, tipo, dataInicio, dataFim)
}

func (service *CrimeService) Update(
	ctx context.Context,
	id bson.ObjectID,
	crime models.Crime,
) error {
	if err := validarCrime(&crime); err != nil {
		return err
	}
	return service.repository.Update(ctx, id, crime)
}

func (service *CrimeService) Delete(ctx context.Context, id bson.ObjectID) error {
	return service.repository.Delete(ctx, id)
}
