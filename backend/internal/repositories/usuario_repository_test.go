package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"AEP-1B-2026-2/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func newUsuarioTestRepository(t *testing.T) (*UsuarioRepository, context.Context) {
	t.Helper()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://aep_root:change_me@localhost:27017/AEP_test?authSource=admin"
	}
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("Mongo indisponível para teste de integração: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("Mongo indisponível para teste de integração: %v", err)
	}
	db := client.Database("AEP_usuario_test_" + time.Now().Format("20060102150405000000000"))
	probe := db.Collection("_test_access")
	if _, err := probe.InsertOne(ctx, bson.M{"test": true}); err != nil {
		_ = client.Disconnect(context.Background())
		t.Skipf("Mongo sem permissão de escrita para teste de integração: %v", err)
	}
	if _, err := probe.DeleteMany(ctx, bson.M{}); err != nil {
		_ = client.Disconnect(context.Background())
		t.Skipf("Mongo sem permissão de limpeza para teste de integração: %v", err)
	}
	t.Cleanup(func() { _ = db.Drop(context.Background()); _ = client.Disconnect(context.Background()) })
	return NewUsuarioRepository(db), ctx
}

func TestUsuarioRepository_Integration(t *testing.T) {
	repository, ctx := newUsuarioTestRepository(t)
	usuario := &models.Usuario{Nome: "Ana", Email: "ana@example.com", Senha: "segredo", Ativo: true}
	if err := repository.Create(ctx, usuario); err != nil {
		t.Fatalf("Create() falhou: %v", err)
	}
	if usuario.ID.IsZero() {
		t.Fatal("Create() não gerou ID")
	}

	usuarios, err := repository.FindAll(ctx)
	if err != nil || len(usuarios) != 1 {
		t.Fatalf("FindAll() = %d usuários, %v; esperado 1 e nil", len(usuarios), err)
	}
	found, err := repository.FindByID(ctx, usuario.ID)
	if err != nil || found.Email != usuario.Email {
		t.Fatalf("FindByID() = %#v, %v", found, err)
	}

	updated := models.Usuario{Nome: "Ana Silva", Email: "ana.silva@example.com", Senha: "nova", Ativo: false}
	if err := repository.Update(ctx, usuario.ID, updated); err != nil {
		t.Fatalf("Update() falhou: %v", err)
	}
	found, err = repository.FindByID(ctx, usuario.ID)
	if err != nil || found.Nome != updated.Nome || found.Ativo != updated.Ativo {
		t.Fatalf("usuário não foi atualizado: %#v, %v", found, err)
	}
	if err := repository.Delete(ctx, usuario.ID); err != nil {
		t.Fatalf("Delete() falhou: %v", err)
	}
	if _, err := repository.FindByID(ctx, usuario.ID); err == nil {
		t.Fatal("FindByID() deveria falhar após Delete()")
	}
}

func TestUsuarioRepository_UpdateAndDeleteMissingDocument(t *testing.T) {
	repository, ctx := newUsuarioTestRepository(t)
	id := bson.NewObjectID()
	if err := repository.Update(ctx, id, models.Usuario{Nome: "Inexistente"}); err != mongo.ErrNoDocuments {
		t.Fatalf("Update() = %v; esperado mongo.ErrNoDocuments", err)
	}
	if err := repository.Delete(ctx, id); err != mongo.ErrNoDocuments {
		t.Fatalf("Delete() = %v; esperado mongo.ErrNoDocuments", err)
	}
}
