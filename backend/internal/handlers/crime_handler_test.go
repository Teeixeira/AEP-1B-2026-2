package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type stubCrimeService struct {
	createFn   func(ctx context.Context, crime *models.Crime) error
	findAllFn  func(ctx context.Context) ([]models.Crime, error)
	findByIDFn func(ctx context.Context, id bson.ObjectID) (*models.Crime, error)
	findNearFn func(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error)
	updateFn   func(ctx context.Context, id bson.ObjectID, crime models.Crime) error
	deleteFn   func(ctx context.Context, id bson.ObjectID) error
}

func (s *stubCrimeService) Create(ctx context.Context, crime *models.Crime) error {
	if s.createFn != nil {
		return s.createFn(ctx, crime)
	}
	return nil
}

func (s *stubCrimeService) FindAll(ctx context.Context) ([]models.Crime, error) {
	if s.findAllFn != nil {
		return s.findAllFn(ctx)
	}
	return nil, nil
}

func (s *stubCrimeService) FindByID(ctx context.Context, id bson.ObjectID) (*models.Crime, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubCrimeService) FindNear(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error) {
	if s.findNearFn != nil {
		return s.findNearFn(ctx, longitude, latitude, raioMetros)
	}
	return nil, nil
}

func (s *stubCrimeService) Update(ctx context.Context, id bson.ObjectID, crime models.Crime) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, crime)
	}
	return nil
}

func (s *stubCrimeService) Delete(ctx context.Context, id bson.ObjectID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

func TestCrimeHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := `{"tipo":"furto","descricao":"aparelho","data_hora":"2026-08-28T12:00:00Z","localizacao":{"type":"Point","coordinates":[-46.6333,-23.5505]}}`
	service := &stubCrimeService{
		createFn: func(ctx context.Context, crime *models.Crime) error {
			if crime == nil {
				t.Fatal("crime não pode ser nulo")
			}
			if crime.Tipo != "furto" {
				t.Fatalf("tipo inesperado: %s", crime.Tipo)
			}
			crime.ID = bson.NewObjectID()
			return nil
		},
	}

	handler := NewCrimeHandler(service)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/crimes", strings.NewReader(payload))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Create(ctx)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
	var resposta models.Crime
	if err := json.Unmarshal(rec.Body.Bytes(), &resposta); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if resposta.Tipo != "furto" {
		t.Fatalf("tipo retornado inesperado: %s", resposta.Tipo)
	}
}

func TestCrimeHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewCrimeHandler(&stubCrimeService{})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/crimes", strings.NewReader(`{"tipo":`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Create(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_FindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	crimesEsperados := []models.Crime{{ID: bson.NewObjectID(), Tipo: "assalto"}}
	handler := NewCrimeHandler(&stubCrimeService{findAllFn: func(ctx context.Context) ([]models.Crime, error) {
		return crimesEsperados, nil
	}})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/crimes", nil)

	handler.FindAll(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_FindByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	handler := NewCrimeHandler(&stubCrimeService{findByIDFn: func(ctx context.Context, idRecebido bson.ObjectID) (*models.Crime, error) {
		if idRecebido != id {
			t.Fatalf("id inesperado: %s", idRecebido.Hex())
		}
		return &models.Crime{ID: id, Tipo: "furto"}, nil
	}})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/crimes/"+id.Hex(), nil)
	ctx.Params = gin.Params{{Key: "id", Value: id.Hex()}}

	handler.FindByID(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_FindByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewCrimeHandler(&stubCrimeService{})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/crimes/id-invalido", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "id-invalido"}}

	handler.FindByID(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_FindByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	handler := NewCrimeHandler(&stubCrimeService{findByIDFn: func(ctx context.Context, idRecebido bson.ObjectID) (*models.Crime, error) {
		return nil, errors.New("not found")
	}})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/crimes/"+id.Hex(), nil)
	ctx.Params = gin.Params{{Key: "id", Value: id.Hex()}}

	handler.FindByID(ctx)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_FindNear(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewCrimeHandler(&stubCrimeService{findNearFn: func(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error) {
		if longitude != -46.63 || latitude != -23.55 || raioMetros != 500 {
			t.Fatalf("valores inesperados: lon=%.2f lat=%.2f raio=%.2f", longitude, latitude, raioMetros)
		}
		return []models.Crime{{Tipo: "roubo"}}, nil
	}})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/crimes/proximos?lat=-23.55&lng=-46.63&raio=500", nil)

	handler.FindNear(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_FindNear_InvalidLatLng(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewCrimeHandler(&stubCrimeService{})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/crimes/proximos?lat=abc&lng=-46.63", nil)

	handler.FindNear(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_FindNear_InvalidRadius(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewCrimeHandler(&stubCrimeService{findNearFn: func(ctx context.Context, longitude float64, latitude float64, raioMetros float64) ([]models.Crime, error) {
		return nil, services.ErrRaioInvalido
	}})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/crimes/proximos?lat=-23.55&lng=-46.63&raio=0", nil)

	handler.FindNear(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	payload := `{"tipo":"assalto","descricao":"atualizado","data_hora":"2026-08-28T15:00:00Z","localizacao":{"type":"Point","coordinates":[-46.60,-23.59]}}`
	handler := NewCrimeHandler(&stubCrimeService{updateFn: func(ctx context.Context, idRecebido bson.ObjectID, crime models.Crime) error {
		if idRecebido != id {
			t.Fatalf("id inesperado: %s", idRecebido.Hex())
		}
		if crime.Tipo != "assalto" {
			t.Fatalf("tipo inesperado: %s", crime.Tipo)
		}
		return nil
	}})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/crimes/"+id.Hex(), strings.NewReader(payload))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "id", Value: id.Hex()}}

	handler.Update(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}

func TestCrimeHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	handler := NewCrimeHandler(&stubCrimeService{deleteFn: func(ctx context.Context, idRecebido bson.ObjectID) error {
		if idRecebido != id {
			t.Fatalf("id inesperado: %s", idRecebido.Hex())
		}
		return nil
	}})
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/crimes/"+id.Hex(), nil)
	ctx.Params = gin.Params{{Key: "id", Value: id.Hex()}}

	handler.Delete(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status inesperado: %d", rec.Code)
	}
}
