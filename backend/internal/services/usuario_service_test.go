package services

import (
	"context"
	"errors"
	"testing"

	"AEP-1B-2026-2/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type stubUsuarioRepository struct {
	createFn   func(context.Context, *models.Usuario) error
	findAllFn  func(context.Context) ([]models.Usuario, error)
	findByIDFn func(context.Context, bson.ObjectID) (*models.Usuario, error)
	updateFn   func(context.Context, bson.ObjectID, models.Usuario) error
	deleteFn   func(context.Context, bson.ObjectID) error
}

func (s *stubUsuarioRepository) Create(c context.Context, u *models.Usuario) error {
	if s.createFn != nil {
		return s.createFn(c, u)
	}
	return nil
}
func (s *stubUsuarioRepository) FindAll(c context.Context) ([]models.Usuario, error) {
	if s.findAllFn != nil {
		return s.findAllFn(c)
	}
	return nil, nil
}
func (s *stubUsuarioRepository) FindByID(c context.Context, id bson.ObjectID) (*models.Usuario, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(c, id)
	}
	return nil, nil
}
func (s *stubUsuarioRepository) Update(c context.Context, id bson.ObjectID, u models.Usuario) error {
	if s.updateFn != nil {
		return s.updateFn(c, id, u)
	}
	return nil
}
func (s *stubUsuarioRepository) Delete(c context.Context, id bson.ObjectID) error {
	if s.deleteFn != nil {
		return s.deleteFn(c, id)
	}
	return nil
}

func TestUsuarioService_DelegatesOperations(t *testing.T) {
	ctx, id, usuario := context.Background(), bson.NewObjectID(), &models.Usuario{Nome: "Ana"}
	repository := &stubUsuarioRepository{
		createFn: func(_ context.Context, got *models.Usuario) error {
			if got != usuario {
				t.Fatal("usuário não repassado")
			}
			return nil
		},
		findAllFn: func(context.Context) ([]models.Usuario, error) { return []models.Usuario{{ID: id}}, nil },
		findByIDFn: func(_ context.Context, got bson.ObjectID) (*models.Usuario, error) {
			if got != id {
				t.Fatal("ID não repassado")
			}
			return usuario, nil
		},
		updateFn: func(_ context.Context, got bson.ObjectID, u models.Usuario) error {
			if got != id || u.Nome != "Atualizada" {
				t.Fatal("dados de atualização incorretos")
			}
			return nil
		},
		deleteFn: func(_ context.Context, got bson.ObjectID) error {
			if got != id {
				t.Fatal("ID não repassado")
			}
			return nil
		},
	}
	service := NewUsuarioService(repository)
	if err := service.Create(ctx, usuario); err != nil {
		t.Fatal(err)
	}
	if users, err := service.FindAll(ctx); err != nil || len(users) != 1 {
		t.Fatalf("FindAll() = %#v, %v", users, err)
	}
	if got, err := service.FindByID(ctx, id); err != nil || got != usuario {
		t.Fatalf("FindByID() = %#v, %v", got, err)
	}
	if err := service.Update(ctx, id, models.Usuario{Nome: "Atualizada"}); err != nil {
		t.Fatal(err)
	}
	if err := service.Delete(ctx, id); err != nil {
		t.Fatal(err)
	}
}

func TestUsuarioService_PropagatesErrors(t *testing.T) {
	errExpected, id := errors.New("falha no repositório"), bson.NewObjectID()
	for name, call := range map[string]func(*UsuarioService) error{
		"Create": func(s *UsuarioService) error { return s.Create(context.Background(), &models.Usuario{}) }, "FindAll": func(s *UsuarioService) error { _, err := s.FindAll(context.Background()); return err }, "FindByID": func(s *UsuarioService) error { _, err := s.FindByID(context.Background(), id); return err }, "Update": func(s *UsuarioService) error { return s.Update(context.Background(), id, models.Usuario{}) }, "Delete": func(s *UsuarioService) error { return s.Delete(context.Background(), id) },
	} {
		t.Run(name, func(t *testing.T) {
			r := &stubUsuarioRepository{createFn: func(context.Context, *models.Usuario) error { return errExpected }, findAllFn: func(context.Context) ([]models.Usuario, error) { return nil, errExpected }, findByIDFn: func(context.Context, bson.ObjectID) (*models.Usuario, error) { return nil, errExpected }, updateFn: func(context.Context, bson.ObjectID, models.Usuario) error { return errExpected }, deleteFn: func(context.Context, bson.ObjectID) error { return errExpected }}
			if err := call(NewUsuarioService(r)); !errors.Is(err, errExpected) {
				t.Fatalf("erro = %v; esperado %v", err, errExpected)
			}
		})
	}
}
