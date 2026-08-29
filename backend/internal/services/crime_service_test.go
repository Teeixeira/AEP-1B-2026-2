package services

import (
	"context"
	"errors"
	"testing"

	"AEP-1B-2026-2/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type stubCrimeRepository struct {
	createFn   func(ctx context.Context, crime *models.Crime) error
	findAllFn  func(ctx context.Context) ([]models.Crime, error)
	findByIDFn func(ctx context.Context, id bson.ObjectID) (*models.Crime, error)
	findNearFn func(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error)
	updateFn   func(ctx context.Context, id bson.ObjectID, crime models.Crime) error
	deleteFn   func(ctx context.Context, id bson.ObjectID) error
	ensureFn   func(ctx context.Context) error
}

func (s *stubCrimeRepository) EnsureGeoIndex(ctx context.Context) error {
	if s.ensureFn != nil {
		return s.ensureFn(ctx)
	}
	return nil
}

func (s *stubCrimeRepository) Create(ctx context.Context, crime *models.Crime) error {
	if s.createFn != nil {
		return s.createFn(ctx, crime)
	}
	return nil
}

func (s *stubCrimeRepository) FindAll(ctx context.Context) ([]models.Crime, error) {
	if s.findAllFn != nil {
		return s.findAllFn(ctx)
	}
	return nil, nil
}

func (s *stubCrimeRepository) FindByID(ctx context.Context, id bson.ObjectID) (*models.Crime, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubCrimeRepository) FindNear(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error) {
	if s.findNearFn != nil {
		return s.findNearFn(ctx, longitude, latitude, raioMetros)
	}
	return nil, nil
}

func (s *stubCrimeRepository) Update(ctx context.Context, id bson.ObjectID, crime models.Crime) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, crime)
	}
	return nil
}

func (s *stubCrimeRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

func TestCrimeService_Create(t *testing.T) {
	crime := &models.Crime{
		Tipo:      "furto",
		Descricao: "carteira roubada",
		DataHora:  "2026-08-28T12:00:00Z",
		Localizacao: models.Localizacao{
			Coordinates: []float64{-46.6333, -23.5505},
		},
	}

	repository := &stubCrimeRepository{
		createFn: func(ctx context.Context, crimeRecebido *models.Crime) error {
			if crimeRecebido == nil {
				t.Fatal("crime não pode ser nulo")
			}
			if crimeRecebido.Localizacao.Type != models.TipoPonto {
				t.Fatalf("tipo de localização inesperado: %s", crimeRecebido.Localizacao.Type)
			}
			return nil
		},
	}

	service := NewCrimeService(repository)
	if err := service.Create(context.Background(), crime); err != nil {
		t.Fatalf("Create() retornou erro: %v", err)
	}
}

func TestCrimeService_FindAll(t *testing.T) {
	crimesEsperados := []models.Crime{{ID: bson.NewObjectID(), Tipo: "furto"}}
	repository := &stubCrimeRepository{
		findAllFn: func(ctx context.Context) ([]models.Crime, error) {
			return crimesEsperados, nil
		},
	}

	service := NewCrimeService(repository)
	result, err := service.FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() retornou erro: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("quantidade inesperada: %d", len(result))
	}
}

func TestCrimeService_FindByID(t *testing.T) {
	id := bson.NewObjectID()
	crimeEsperada := &models.Crime{ID: id, Tipo: "assalto"}
	repository := &stubCrimeRepository{
		findByIDFn: func(ctx context.Context, idRecebido bson.ObjectID) (*models.Crime, error) {
			if idRecebido != id {
				t.Fatalf("id inesperado: %s", idRecebido.Hex())
			}
			return crimeEsperada, nil
		},
	}

	service := NewCrimeService(repository)
	result, err := service.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID() retornou erro: %v", err)
	}
	if result == nil || result.ID != id {
		t.Fatalf("resultado inesperado: %#v", result)
	}
}

func TestCrimeService_FindNear(t *testing.T) {
	crimesEsperados := []models.Crime{{Tipo: "roubo"}}
	repository := &stubCrimeRepository{
		findNearFn: func(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error) {
			if longitude != -46.63 || latitude != -23.55 || raioMetros != 500 {
				t.Fatalf("parametros inesperados: lon=%.2f lat=%.2f raio=%.2f", longitude, latitude, raioMetros)
			}
			return crimesEsperados, nil
		},
	}

	service := NewCrimeService(repository)
	result, err := service.FindNear(context.Background(), -46.63, -23.55, 500)
	if err != nil {
		t.Fatalf("FindNear() retornou erro: %v", err)
	}
	if len(result) != len(crimesEsperados) {
		t.Fatalf("quantidade inesperada: %d", len(result))
	}
}

func TestCrimeService_FindNear_ValidaRaio(t *testing.T) {
	service := NewCrimeService(&stubCrimeRepository{})
	result, err := service.FindNear(context.Background(), -46.63, -23.55, 0)
	if !errors.Is(err, ErrRaioInvalido) {
		t.Fatalf("erro esperado: %v, recebido: %v", ErrRaioInvalido, err)
	}
	if result != nil {
		t.Fatalf("resultado esperado nulo, recebido: %#v", result)
	}
}

func TestCrimeService_Update(t *testing.T) {
	id := bson.NewObjectID()
	crimeAtualizada := models.Crime{
		Tipo:      "furto",
		Descricao: "atualizado",
		DataHora:  "2026-08-28T13:30:00Z",
		Localizacao: models.Localizacao{
			Coordinates: []float64{-46.63, -23.54},
		},
	}

	repository := &stubCrimeRepository{
		updateFn: func(ctx context.Context, idRecebido bson.ObjectID, crime models.Crime) error {
			if idRecebido != id {
				t.Fatalf("id inesperado: %s", idRecebido.Hex())
			}
			if crime.Localizacao.Type != models.TipoPonto {
				t.Fatalf("tipo de localização inesperado: %s", crime.Localizacao.Type)
			}
			return nil
		},
	}

	service := NewCrimeService(repository)
	if err := service.Update(context.Background(), id, crimeAtualizada); err != nil {
		t.Fatalf("Update() retornou erro: %v", err)
	}
}

func TestCrimeService_Delete(t *testing.T) {
	id := bson.NewObjectID()
	repository := &stubCrimeRepository{
		deleteFn: func(ctx context.Context, idRecebido bson.ObjectID) error {
			if idRecebido != id {
				t.Fatalf("id inesperado: %s", idRecebido.Hex())
			}
			return nil
		},
	}

	service := NewCrimeService(repository)
	if err := service.Delete(context.Background(), id); err != nil {
		t.Fatalf("Delete() retornou erro: %v", err)
	}
}

func TestCrimeService_PropagaErros(t *testing.T) {
	expectedErr := errors.New("falha de persistência")
	id := bson.NewObjectID()

	t.Run("Create", func(t *testing.T) {
		repository := &stubCrimeRepository{createFn: func(ctx context.Context, crime *models.Crime) error { return expectedErr }}
		service := NewCrimeService(repository)
		if err := service.Create(context.Background(), &models.Crime{Tipo: "furto"}); !errors.Is(err, expectedErr) {
			t.Fatalf("erro esperado %v, recebido %v", expectedErr, err)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		repository := &stubCrimeRepository{findAllFn: func(ctx context.Context) ([]models.Crime, error) { return nil, expectedErr }}
		service := NewCrimeService(repository)
		if _, err := service.FindAll(context.Background()); !errors.Is(err, expectedErr) {
			t.Fatalf("erro esperado %v, recebido %v", expectedErr, err)
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		repository := &stubCrimeRepository{findByIDFn: func(ctx context.Context, id bson.ObjectID) (*models.Crime, error) { return nil, expectedErr }}
		service := NewCrimeService(repository)
		if _, err := service.FindByID(context.Background(), id); !errors.Is(err, expectedErr) {
			t.Fatalf("erro esperado %v, recebido %v", expectedErr, err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		repository := &stubCrimeRepository{updateFn: func(ctx context.Context, id bson.ObjectID, crime models.Crime) error { return expectedErr }}
		service := NewCrimeService(repository)
		if err := service.Update(context.Background(), id, models.Crime{Tipo: "roubo"}); !errors.Is(err, expectedErr) {
			t.Fatalf("erro esperado %v, recebido %v", expectedErr, err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		repository := &stubCrimeRepository{deleteFn: func(ctx context.Context, id bson.ObjectID) error { return expectedErr }}
		service := NewCrimeService(repository)
		if err := service.Delete(context.Background(), id); !errors.Is(err, expectedErr) {
			t.Fatalf("erro esperado %v, recebido %v", expectedErr, err)
		}
	})
}
