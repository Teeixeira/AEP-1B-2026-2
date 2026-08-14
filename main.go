package main

import (
	"fmt"
	"log"

	"AEP-1B-2026-2/internal/database"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}

	client, err := database.Connect()

	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}

	defer client.Disconnect(nil)

	fmt.Println("MongoDB conectado com sucesso!")
}
