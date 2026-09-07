package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"AEP-1B-2026-2/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type stubUsuarioService struct {
	createFn   func(context.Context, *models.Usuario) error
	findAllFn  func(context.Context) ([]models.Usuario, error)
	findByIDFn func(context.Context, bson.ObjectID) (*models.Usuario, error)
	updateFn   func(context.Context, bson.ObjectID, models.Usuario) error
	deleteFn   func(context.Context, bson.ObjectID) error
}

func (s *stubUsuarioService) Create(c context.Context, u *models.Usuario) error {
	if s.createFn != nil {
		return s.createFn(c, u)
	}
	return nil
}
func (s *stubUsuarioService) FindAll(c context.Context) ([]models.Usuario, error) {
	if s.findAllFn != nil {
		return s.findAllFn(c)
	}
	return nil, nil
}
func (s *stubUsuarioService) FindByID(c context.Context, id bson.ObjectID) (*models.Usuario, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(c, id)
	}
	return nil, nil
}
func (s *stubUsuarioService) Update(c context.Context, id bson.ObjectID, u models.Usuario) error {
	if s.updateFn != nil {
		return s.updateFn(c, id, u)
	}
	return nil
}
func (s *stubUsuarioService) Delete(c context.Context, id bson.ObjectID) error {
	if s.deleteFn != nil {
		return s.deleteFn(c, id)
	}
	return nil
}

func usuarioContext(method, path, body, id string) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	if id != "" {
		ctx.Params = gin.Params{{Key: "id", Value: id}}
	}
	return ctx, rec
}

func TestUsuarioHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewUsuarioHandler(&stubUsuarioService{createFn: func(_ context.Context, u *models.Usuario) error {
		if u.Nome != "Ana" {
			t.Fatal("JSON não foi convertido")
		}
		u.ID = bson.NewObjectID()
		return nil
	}})
	ctx, rec := usuarioContext(http.MethodPost, "/usuarios", `{"nome":"Ana","email":"ana@example.com","senha":"x","ativo":true}`, "")
	handler.Create(ctx)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestUsuarioHandler_CreateErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for name, tc := range map[string]struct {
		handler  *UsuarioHandler
		body     string
		expected int
	}{
		"invalid JSON":  {NewUsuarioHandler(&stubUsuarioService{}), `{"nome":`, http.StatusBadRequest},
		"service error": {NewUsuarioHandler(&stubUsuarioService{createFn: func(context.Context, *models.Usuario) error { return errors.New("erro") }}), `{}`, http.StatusInternalServerError},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, rec := usuarioContext(http.MethodPost, "/usuarios", tc.body, "")
			tc.handler.Create(ctx)
			if rec.Code != tc.expected {
				t.Fatalf("status = %d", rec.Code)
			}
		})
	}
}

func TestUsuarioHandler_FindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for name, tc := range map[string]struct {
		handler  *UsuarioHandler
		expected int
	}{
		"success": {NewUsuarioHandler(&stubUsuarioService{findAllFn: func(context.Context) ([]models.Usuario, error) { return []models.Usuario{{Nome: "Ana"}}, nil }}), http.StatusOK},
		"error":   {NewUsuarioHandler(&stubUsuarioService{findAllFn: func(context.Context) ([]models.Usuario, error) { return nil, errors.New("erro") }}), http.StatusInternalServerError},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, rec := usuarioContext(http.MethodGet, "/usuarios", "", "")
			tc.handler.FindAll(ctx)
			if rec.Code != tc.expected {
				t.Fatalf("status = %d", rec.Code)
			}
		})
	}
}

func TestUsuarioHandler_FindByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	for name, tc := range map[string]struct {
		handler  *UsuarioHandler
		value    string
		expected int
	}{
		"success": {NewUsuarioHandler(&stubUsuarioService{findByIDFn: func(_ context.Context, got bson.ObjectID) (*models.Usuario, error) {
			if got != id {
				t.Fatal("ID incorreto")
			}
			return &models.Usuario{ID: id}, nil
		}}), id.Hex(), http.StatusOK},
		"invalid ID": {NewUsuarioHandler(&stubUsuarioService{}), "invalido", http.StatusBadRequest},
		"not found": {NewUsuarioHandler(&stubUsuarioService{findByIDFn: func(context.Context, bson.ObjectID) (*models.Usuario, error) {
			return nil, errors.New("não encontrado")
		}}), id.Hex(), http.StatusNotFound},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, rec := usuarioContext(http.MethodGet, "/usuarios/"+tc.value, "", tc.value)
			tc.handler.FindByID(ctx)
			if rec.Code != tc.expected {
				t.Fatalf("status = %d", rec.Code)
			}
		})
	}
}

func TestUsuarioHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	body := `{"nome":"Ana","email":"ana@example.com"}`
	for name, tc := range map[string]struct {
		handler        *UsuarioHandler
		value, payload string
		expected       int
	}{
		"success": {NewUsuarioHandler(&stubUsuarioService{updateFn: func(_ context.Context, got bson.ObjectID, u models.Usuario) error {
			if got != id || u.Nome != "Ana" {
				t.Fatal("dados incorretos")
			}
			return nil
		}}), id.Hex(), body, http.StatusOK},
		"invalid ID":   {NewUsuarioHandler(&stubUsuarioService{}), "invalido", body, http.StatusBadRequest},
		"invalid JSON": {NewUsuarioHandler(&stubUsuarioService{}), id.Hex(), `{`, http.StatusBadRequest},
		"not found":    {NewUsuarioHandler(&stubUsuarioService{updateFn: func(context.Context, bson.ObjectID, models.Usuario) error { return errors.New("erro") }}), id.Hex(), body, http.StatusNotFound},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, rec := usuarioContext(http.MethodPut, "/usuarios/"+tc.value, tc.payload, tc.value)
			tc.handler.Update(ctx)
			if rec.Code != tc.expected {
				t.Fatalf("status = %d", rec.Code)
			}
		})
	}
}

func TestUsuarioHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := bson.NewObjectID()
	for name, tc := range map[string]struct {
		handler  *UsuarioHandler
		value    string
		expected int
	}{
		"success": {NewUsuarioHandler(&stubUsuarioService{deleteFn: func(_ context.Context, got bson.ObjectID) error {
			if got != id {
				t.Fatal("ID incorreto")
			}
			return nil
		}}), id.Hex(), http.StatusOK},
		"invalid ID": {NewUsuarioHandler(&stubUsuarioService{}), "invalido", http.StatusBadRequest},
		"not found":  {NewUsuarioHandler(&stubUsuarioService{deleteFn: func(context.Context, bson.ObjectID) error { return errors.New("erro") }}), id.Hex(), http.StatusNotFound},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, rec := usuarioContext(http.MethodDelete, "/usuarios/"+tc.value, "", tc.value)
			tc.handler.Delete(ctx)
			if rec.Code != tc.expected {
				t.Fatalf("status = %d", rec.Code)
			}
		})
	}
}
