package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"AEP-1B-2026-2/internal/database"
	"AEP-1B-2026-2/internal/handlers"
	_ "AEP-1B-2026-2/internal/migrations"
	"AEP-1B-2026-2/internal/repositories"
	"AEP-1B-2026-2/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	migrate "github.com/xakep666/mongo-migrate"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "AEP-1B-2026-2/docs"
)

// @title AEP API
// @version 1.0
// @description API do projeto AEP para consulta e gerenciamento de ocorrências criminais e usuários.
// @BasePath /api
func main() {
	migrateDown := flag.Bool("migrate-down", false, "Reverte a última migração aplicada e encerra")
	migrateOnly := flag.Bool("migrate-only", false, "Executa as migrações pendentes e encerra sem iniciar o servidor")
	flag.Parse()

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

	migrate.SetDatabase(databaseMongo)

	// Se solicitado rollback via flag:
	if *migrateDown {
		log.Println("Executando rollback da última migração...")
		if err := migrate.Down(context.Background(), 1); err != nil {
			log.Fatal("Erro ao reverter migração:", err)
		}
		log.Println("Rollback concluído com sucesso!")
		return
	}

	// Executa migrações pendentes no banco
	if err := migrate.Up(context.Background(), migrate.AllAvailable); err != nil {
		log.Fatal("Erro ao executar migrações:", err)
	}
	log.Println("Migrações executadas com sucesso!")

	// Se solicitado apenas migrar e sair:
	if *migrateOnly {
		return
	}

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
	log.Println("Swagger path: /swagger/index.html")

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
