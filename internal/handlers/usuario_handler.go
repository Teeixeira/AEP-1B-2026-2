package handlers

import (
	"net/http"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UsuarioHandler struct {
	service *services.UsuarioService
}

func NewUsuarioHandler(
	service *services.UsuarioService,
) *UsuarioHandler {

	return &UsuarioHandler{
		service: service,
	}
}

func (handler *UsuarioHandler) Create(c *gin.Context) {

	var usuario models.Usuario

	err := c.ShouldBindJSON(&usuario)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON Inválido",
		})
		return
	}

	err = handler.service.Create(
		c.Request.Context(),
		&usuario,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao criar usuário",
		})
		return
	}

	c.JSON(http.StatusCreated, usuario)
}

func (handler *UsuarioHandler) FindAll(c *gin.Context) {

	usuarios, err := handler.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao buscar usuários",
		})
		return
	}

	c.JSON(http.StatusOK, usuarios)
}

func (handler *UsuarioHandler) FindByID(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID Inválido",
		})
		return
	}

	usuario, err := handler.service.FindByID(
		c.Request.Context(),
		objectID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Usuário não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

func (handler *UsuarioHandler) Update(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID Inválido",
		})
		return
	}

	var usuario models.Usuario

	err = c.ShouldBindJSON(&usuario)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON Inválido",
		})
		return
	}

	err = handler.service.Update(
		c.Request.Context(),
		objectID,
		usuario,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Usuário não encontrado",
		})
		return
	}

	usuario.ID = objectID

	c.JSON(http.StatusOK, usuario)
}

func (handler *UsuarioHandler) Delete(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	err = handler.service.Delete(
		c.Request.Context(),
		objectID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "Usuário não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Usuário excluído com sucesso",
	})
}
