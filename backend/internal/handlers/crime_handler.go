package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CrimeHandler struct {
	service services.CrimeServiceInterface
}

func NewCrimeHandler(service services.CrimeServiceInterface) *CrimeHandler {
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
// @Failure 400 {object} ErroResposta
// @Failure 500 {object} ErroResposta
// @Router /crimes [post]
func (handler *CrimeHandler) Create(c *gin.Context) {
	var crime models.Crime
	if err := c.ShouldBindJSON(&crime); err != nil {
		RespostaEntradaInvalida(c, "JSON Inválido", err.Error())
		return
	}

	if err := handler.service.Create(c.Request.Context(), &crime); err != nil {
		if errors.Is(err, services.ErrCamposObrigatorios) ||
			errors.Is(err, services.ErrLatitudeInvalida) ||
			errors.Is(err, services.ErrLongitudeInvalida) {
			RespostaEntradaInvalida(c, err.Error())
			return
		}
		RespostaFalhaInterna(c, "Erro ao criar crime", err.Error())
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
// @Failure 500 {object} ErroResposta
// @Router /crimes [get]
func (handler *CrimeHandler) FindAll(c *gin.Context) {
	crimes, err := handler.service.FindAll(c.Request.Context())
	if err != nil {
		RespostaFalhaInterna(c, "Erro ao buscar crimes", err.Error())
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
// @Failure 400 {object} ErroResposta
// @Failure 404 {object} ErroResposta
// @Router /crimes/{id} [get]
func (handler *CrimeHandler) FindByID(c *gin.Context) {
	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespostaEntradaInvalida(c, "ID Inválido")
		return
	}

	crime, err := handler.service.FindByID(c.Request.Context(), objectID)
	if err != nil {
		RespostaNaoEncontrado(c, "Crime não encontrado")
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
// @Param raio query number false "Raio de busca em metros" default(1000)
// @Param tipo query string false "Tipo de crime"
// @Param data_inicio query string false "Início do período, em RFC3339"
// @Param data_fim query string false "Fim do período, em RFC3339"
// @Success 200 {array} models.Crime
// @Failure 400 {object} ErroResposta
// @Failure 500 {object} ErroResposta
// @Router /crimes/proximos [get]
func (handler *CrimeHandler) FindNear(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		RespostaEntradaInvalida(c, "Parâmetro 'lat' inválido")
		return
	}

	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		RespostaEntradaInvalida(c, "Parâmetro 'lng' inválido")
		return
	}

	raio, err := strconv.ParseFloat(c.DefaultQuery("raio", "1000"), 64)
	if err != nil {
		RespostaEntradaInvalida(c, "Parâmetro 'raio' inválido")
		return
	}

	tipo := c.Query("tipo")
	dataInicio, err := parseOptionalDateTime(c.Query("data_inicio"))
	if err != nil {
		RespostaEntradaInvalida(c, "Parâmetro 'data_inicio' inválido; use RFC3339")
		return
	}
	dataFim, err := parseOptionalDateTime(c.Query("data_fim"))
	if err != nil {
		RespostaEntradaInvalida(c, "Parâmetro 'data_fim' inválido; use RFC3339")
		return
	}

	crimes, err := handler.service.FindNear(
		c.Request.Context(),
		lng,
		lat,
		raio,
		tipo,
		dataInicio,
		dataFim,
	)

	if err != nil {
		if errors.Is(err, services.ErrRaioInvalido) ||
			errors.Is(err, services.ErrPeriodoInvalido) {
			RespostaEntradaInvalida(c, err.Error())
			return
		}

		RespostaFalhaInterna(c, "Erro ao buscar crimes próximos", err.Error())
		return
	}

	c.JSON(http.StatusOK, crimes)
}

func parseOptionalDateTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, value)
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
// @Failure 400 {object} ErroResposta
// @Failure 404 {object} ErroResposta
// @Failure 500 {object} ErroResposta
// @Router /crimes/{id} [put]
func (handler *CrimeHandler) Update(c *gin.Context) {
	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespostaEntradaInvalida(c, "ID Inválido")
		return
	}

	var crime models.Crime
	if err := c.ShouldBindJSON(&crime); err != nil {
		RespostaEntradaInvalida(c, "JSON Inválido", err.Error())
		return
	}

	err = handler.service.Update(
		c.Request.Context(),
		objectID,
		crime,
	)

	if err != nil {
		if errors.Is(err, services.ErrCamposObrigatorios) ||
			errors.Is(err, services.ErrLatitudeInvalida) ||
			errors.Is(err, services.ErrLongitudeInvalida) {
			RespostaEntradaInvalida(c, err.Error())
			return
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			RespostaNaoEncontrado(c, "Crime não encontrado")
			return
		}

		RespostaFalhaInterna(c, "Erro ao atualizar crime", err.Error())
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
// @Failure 400 {object} ErroResposta
// @Failure 404 {object} ErroResposta
// @Failure 500 {object} ErroResposta
// @Router /crimes/{id} [delete]
func (handler *CrimeHandler) Delete(c *gin.Context) {
	objectID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespostaEntradaInvalida(c, "ID inválido")
		return
	}

	err = handler.service.Delete(
		c.Request.Context(),
		objectID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			RespostaNaoEncontrado(c, "Crime não encontrado")
			return
		}

		RespostaFalhaInterna(c, "Erro ao excluir crime", err.Error())
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"mensagem": "Crime excluído com sucesso"},
	)
}
