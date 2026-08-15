package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"AEP-1B-2026-2/internal/database"
	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/repositories"

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

	databaseMongo := client.Database("AEP")

	repository := repositories.NewUsuarioRepository(databaseMongo)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	usuario := models.Usuario{
		Nome:  "Leonardo",
		Email: "leonardo@email.com",
		Senha: "123456",
		Ativo: true,
	}

	err = repository.Create(ctx, &usuario)

	if err != nil {
		log.Fatal("Erro ao criar usuário:", err)
	}

	fmt.Println("Usuário criado com suscesso!")
	fmt.Println("ID:", usuario.ID)

	usuarios, err := repository.FindAll(ctx)

	if err != nil {
		log.Fatal("Erro ao buscar usuários:", err)
	}

	fmt.Println("Usuários encontrados:")

	for _, usuario := range usuarios {
		fmt.Println(usuario)
	}

	usuarioEncontrado, err := repository.FindByID(
		ctx,
		usuario.ID,
	)

	if err != nil {
		log.Fatal("Erro ao buscar usuário:", err)
	}

	fmt.Println("Usuário encontrado:")
	fmt.Println("ID:", usuarioEncontrado.ID)
	fmt.Println("Nome:", usuarioEncontrado.Nome)
	fmt.Println("Email:", usuarioEncontrado.Email)
	fmt.Println("Ativo:", usuarioEncontrado.Ativo)
}
