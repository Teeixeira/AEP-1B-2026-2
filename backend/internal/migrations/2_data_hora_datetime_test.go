package migrations

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestParseDataHora(t *testing.T) {
	value := "2026-08-28T15:00:00.123+03:00"
	got, err := parseDataHora(value)
	if err != nil {
		t.Fatalf("parseDataHora() retornou erro: %v", err)
	}

	want, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatalf("falha ao preparar data esperada: %v", err)
	}
	if !got.Equal(want) {
		t.Fatalf("data inesperada: recebido %v, esperado %v", got, want)
	}
}

func TestParseDataHora_Invalid(t *testing.T) {
	if _, err := parseDataHora("not-a-date"); err == nil {
		t.Fatal("parseDataHora() deveria rejeitar uma data inválida")
	}
}

func TestParseDataHora_WithoutTimezoneUsesUTC(t *testing.T) {
	got, err := parseDataHora("2026-08-01T08:30:00")
	if err != nil {
		t.Fatalf("parseDataHora() retornou erro: %v", err)
	}

	want := time.Date(2026, time.August, 1, 8, 30, 0, 0, time.UTC)
	if !got.Equal(want) || got.Location() != time.UTC {
		t.Fatalf("data inesperada: recebido %v, esperado %v", got, want)
	}
}

func TestBSONDateDecodesAsTime(t *testing.T) {
	want := time.Date(2026, time.August, 28, 12, 0, 0, 0, time.UTC)
	encoded, err := bson.Marshal(bson.M{"data_hora": want})
	if err != nil {
		t.Fatalf("bson.Marshal() retornou erro: %v", err)
	}

	var doc struct {
		DataHora time.Time `bson:"data_hora"`
	}
	if err := bson.Unmarshal(encoded, &doc); err != nil {
		t.Fatalf("bson.Unmarshal() retornou erro: %v", err)
	}
	if !doc.DataHora.Equal(want) {
		t.Fatalf("data BSON inesperada: %v", doc.DataHora)
	}
}
