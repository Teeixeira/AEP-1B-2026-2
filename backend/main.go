package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"AEP-1B-2026-2/internal/database"
	"AEP-1B-2026-2/internal/handlers"
	"AEP-1B-2026-2/internal/repositories"
	"AEP-1B-2026-2/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

	_ "AEP-1B-2026-2/docs"
)

// @title AEP API
// @version 1.0
// @description API for the AEP project.
// @host localhost:8080
// @BasePath /api
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado; usando variáveis de ambiente")
	}

	client, err := database.Connect()

	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}

	defer client.Disconnect(nil)

	log.Println("MongoDB conectado com sucesso!")

	databaseMongo := client.Database("AEP")

	router := gin.Default()
	
    swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

    router.GET("/swagger/*any", func(c *gin.Context) {
        if c.Param("any") == "/" {
            c.Redirect(http.StatusFound, "/swagger/index.html")
            return
        }

        swaggerHandler(c)
    })
	
	api := router.Group("/api")

	usuarioRepository := repositories.NewUsuarioRepository(databaseMongo)
	usuarioService := services.NewUsuarioService(usuarioRepository)
	usuarioHandler := handlers.NewUsuarioHandler(usuarioService)

	api.POST("/usuarios", usuarioHandler.Create)
	api.GET("/usuarios", usuarioHandler.FindAll)
	api.GET("/usuarios/:id", usuarioHandler.FindByID)
	api.PUT("/usuarios/:id", usuarioHandler.Update)
	api.DELETE("/usuarios/:id", usuarioHandler.Delete)

	crimeRepository := repositories.NewCrimeRepository(databaseMongo)

	indexCtx, cancelIndex := context.WithTimeout(context.Background(), 10*time.Second)
	err = crimeRepository.EnsureGeoIndex(indexCtx)
	cancelIndex()
	if err != nil {
		log.Fatal("Erro ao criar índice geoespacial:", err)
	}

	crimeService := services.NewCrimeService(crimeRepository)
	crimeHandler := handlers.NewCrimeHandler(crimeService)

	api.POST("/crimes", crimeHandler.Create)
	api.GET("/crimes", crimeHandler.FindAll)
	api.GET("/crimes/proximos", crimeHandler.FindNear)
	api.GET("/crimes/:id", crimeHandler.FindByID)
	api.PUT("/crimes/:id", crimeHandler.Update)
	api.DELETE("/crimes/:id", crimeHandler.Delete)

	log.Println("Servidor rodando na porta 8080")

	err = router.Run(":8080")

	if err != nil {
		log.Fatal(err)
	}

	// ctx, cancel := context.WithTimeout(
	// 	context.Background(),
	// 	10*time.Second,
	// )
	// defer cancel()

	// usuario := models.Usuario{
	// 	Nome:  "Leonardo",
	// 	Email: "leonardo@email.com",
	// 	Senha: "123456",
	// 	Ativo: true,
	// }

	// err = repository.Create(ctx, &usuario)

	// if err != nil {
	// 	log.Fatal("Erro ao criar usuário:", err)
	// }

	// fmt.Println("Usuário criado com suscesso!")
	// fmt.Println("ID:", usuario.ID)

	// usuarios, err := repository.FindAll(ctx)

	// if err != nil {
	// 	log.Fatal("Erro ao buscar usuários:", err)
	// }

	// fmt.Println("Usuários encontrados:")

	// for _, usuario := range usuarios {
	// 	fmt.Println(usuario)
	// }

	// usuarioEncontrado, err := repository.FindByID(
	// 	ctx,
	// 	usuario.ID,
	// )

	// if err != nil {
	// 	log.Fatal("Erro ao buscar usuário:", err)
	// }

	// fmt.Println("Usuário encontrado:")
	// fmt.Println("ID:", usuarioEncontrado.ID)
	// fmt.Println("Nome:", usuarioEncontrado.Nome)
	// fmt.Println("Email:", usuarioEncontrado.Email)
	// fmt.Println("Ativo:", usuarioEncontrado.Ativo)

	// fmt.Println("=============== UPDATE TEST ===============")

	// usuarioAtualizado := models.Usuario{
	// 	Nome:  "Leonardo Atualizado",
	// 	Email: "leonardo.novo@email.com",
	// 	Senha: "654321",
	// 	Ativo: true,
	// }

	// err = repository.Update(
	// 	ctx,
	// 	usuario.ID,
	// 	usuarioAtualizado,
	// )

	// if err != nil {
	// 	log.Fatal("Erro ao atualizar usuário:", err)
	// }

	// fmt.Println("Usuário atualizado com sucesso!")

	// fmt.Println("=============== DELETE TEST ===============")

	// err = repository.Delete(
	// 	ctx,
	// 	usuario.ID,
	// )

	// if err != nil {
	// 	log.Fatal("Erro ao excluir usuário!", err)
	// }

	// fmt.Println("Usuário deletado com sucesso!")

	// fmt.Println("=============== DELETE CONFIRMATION TEST ===============")

	// _, err = repository.FindByID(
	// 	ctx,
	// 	usuario.ID,
	// )

	// if err != nil {
	// 	fmt.Println("Usuário não encontrado após a exclusão!", err)
	// }
}
