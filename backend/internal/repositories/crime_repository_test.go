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

func TestCrimeRepository_Integration(t *testing.T) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017/AEP_test"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("Mongo não disponível para teste de integração: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("Mongo indisponível no ambiente: %v", err)
	}

	databaseName := "AEP_crime_test_" + time.Now().Format("20060102150405")
	db := client.Database(databaseName)
	repository := NewCrimeRepository(db)

	if err := repository.EnsureGeoIndex(ctx); err != nil {
		t.Fatalf("EnsureGeoIndex() falhou: %v", err)
	}

	crime := &models.Crime{
		Tipo:      "furto",
		Descricao: "carteira roubada",
		DataHora:  "2026-08-28T13:00:00Z",
		Localizacao: models.Localizacao{
			Type:        models.TipoPonto,
			Coordinates: []float64{-46.6333, -23.5505},
		},
	}

	if err := repository.Create(ctx, crime); err != nil {
		t.Fatalf("Create() falhou: %v", err)
	}
	if crime.ID.IsZero() {
		t.Fatal("ID não foi gerado")
	}

	list, err := repository.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() falhou: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("FindAll() não retornou crimes")
	}

	found, err := repository.FindByID(ctx, crime.ID)
	if err != nil {
		t.Fatalf("FindByID() falhou: %v", err)
	}
	if found == nil || found.ID != crime.ID {
		t.Fatal("crime encontrado não corresponde ao ID esperado")
	}

	near, err := repository.FindNear(ctx, -46.6333, -23.5505, 1000)
	if err != nil {
		t.Fatalf("FindNear() falhou: %v", err)
	}
	if len(near) == 0 {
		t.Fatal("FindNear() não encontrou crime próximo")
	}

	crimeAtualizada := models.Crime{
		Tipo:      "assalto",
		Descricao: "atualizado",
		DataHora:  "2026-08-28T14:00:00Z",
		Localizacao: models.Localizacao{
			Type:        models.TipoPonto,
			Coordinates: []float64{-46.7000, -23.6000},
		},
	}
	if err := repository.Update(ctx, crime.ID, crimeAtualizada); err != nil {
		t.Fatalf("Update() falhou: %v", err)
	}

	if err := repository.Delete(ctx, crime.ID); err != nil {
		t.Fatalf("Delete() falhou: %v", err)
	}

	if _, err := repository.FindByID(ctx, crime.ID); err == nil {
		t.Fatal("FindByID() deveria falhar depois do delete")
	}

	if err := db.Drop(ctx); err != nil {
		t.Fatalf("Drop() falhou: %v", err)
	}
	if err := client.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect() falhou: %v", err)
	}
}

func TestCrimeRepository_UpdateMissingDocument(t *testing.T) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017/AEP_test"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("Mongo não disponível para teste de integração: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("Mongo indisponível no ambiente: %v", err)
	}

	db := client.Database("AEP_crime_missing_test_" + time.Now().Format("20060102150405"))
	repository := NewCrimeRepository(db)

	id := bson.NewObjectID()
	err = repository.Update(ctx, id, models.Crime{Tipo: "furto"})
	if err == nil {
		t.Fatal("Update() deveria falhar para documento inexistente")
	}
	if err := client.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect() falhou: %v", err)
	}
}
