package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CrimeHandler struct {
	service *services.CrimeService
}

func NewCrimeHandler(service *services.CrimeService) *CrimeHandler {
	return &CrimeHandler{service: service}
}

// Create godoc
// @Summary Criar crime
// @Description Cria um novo registro de crime
// @Tags crimes
// @Accept json
// @Produce json
// @Param crime body models.Crime true "Dados do crime"
// @Success 201 {object} models.Crime
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /crimes [post]
func (handler *CrimeHandler) Create(c *gin.Context) {
	var crime models.Crime
	if err := c.ShouldBindJSON(&crime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON Inválido", "details": err.Error()})
		return
	}

	if err := handler.service.Create(c.Request.Context(), &crime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao criar crime"})
		return
	}

	c.JSON(http.StatusCreated, crime)
}

// FindAll godoc
// @Summary Listar crimes
// @Description Retorna todos os registros de crimes
// @Tags crimes
// @Produce json
// @Success 200 {array} models.Crime
// @Failure 500 {object} map[string]string
// @Router /crimes [get]
func (handler *CrimeHandler) FindAll(c *gin.Context) {
	crimes, err := handler.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar crimes"})
		return
	}

	c.JSON(http.StatusOK, crimes)
}

// FindByID godoc
// @Summary Buscar crime por ID
// @Description Retorna um crime específico pelo seu ID
// @Tags crimes
// @Produce json
// @Param id path string true "ID do crime"
// @Success 200 {object} models.Crime
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /crimes/{id} [get]
func (handler *CrimeHandler) FindByID(c *gin.Context) {
	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID Inválido"})
		return
	}

	crime, err := handler.service.FindByID(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Crime não encontrado"})
		return
	}

	c.JSON(http.StatusOK, crime)
}

// FindNear godoc
// @Summary Buscar crimes próximos
// @Description Retorna crimes próximos a uma localização geográfica
// @Tags crimes
// @Produce json
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Param raio query number false "Raio em metros" default(1000)
// @Success 200 {array} models.Crime
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /crimes/proximos [get]
func (handler *CrimeHandler) FindNear(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Parâmetro 'lat' inválido"})
		return
	}

	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Parâmetro 'lng' inválido"})
		return
	}

	raio, err := strconv.ParseFloat(c.DefaultQuery("raio", "1000"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Parâmetro 'raio' inválido"})
		return
	}

	crimes, err := handler.service.FindNear(
		c.Request.Context(),
		lng,
		lat,
		raio,
	)

	if err != nil {
		if errors.Is(err, services.ErrRaioInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{"erro": "Erro ao buscar crimes próximos"},
		)
		return
	}

	c.JSON(http.StatusOK, crimes)
}

// Update godoc
// @Summary Atualizar crime
// @Description Atualiza um crime existente
// @Tags crimes
// @Accept json
// @Produce json
// @Param id path string true "ID do crime"
// @Param crime body models.Crime true "Dados atualizados do crime"
// @Success 200 {object} models.Crime
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /crimes/{id} [put]
func (handler *CrimeHandler) Update(c *gin.Context) {
	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID Inválido"})
		return
	}

	var crime models.Crime
	if err := c.ShouldBindJSON(&crime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON Inválido"})
		return
	}

	err = handler.service.Update(
		c.Request.Context(),
		objectID,
		crime,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"erro": "Crime não encontrado"})
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{"erro": "Erro ao atualizar crime"},
		)
		return
	}

	crime.ID = objectID
	c.JSON(http.StatusOK, crime)
}

// Delete godoc
// @Summary Excluir crime
// @Description Exclui um crime existente
// @Tags crimes
// @Produce json
// @Param id path string true "ID do crime"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /crimes/{id} [delete]
func (handler *CrimeHandler) Delete(c *gin.Context) {
	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	err = handler.service.Delete(
		c.Request.Context(),
		objectID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"erro": "Crime não encontrado"})
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{"erro": "Erro ao excluir crime"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"mensagem": "Crime excluído com sucesso"},
	)
}
