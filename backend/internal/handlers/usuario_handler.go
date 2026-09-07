package handlers

import (
	"net/http"

	"AEP-1B-2026-2/internal/models"
	"AEP-1B-2026-2/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UsuarioHandler struct {
	service services.UsuarioServiceInterface
}

func NewUsuarioHandler(
	service services.UsuarioServiceInterface,
) *UsuarioHandler {

	return &UsuarioHandler{
		service: service,
	}
}

// Create godoc
// @Summary Criar usuário
// @Description Cria um novo usuário
// @Tags usuarios
// @Accept json
// @Produce json
// @Param usuario body models.Usuario true "Dados do usuário"
// @Success 201 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /usuarios [post]
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

// FindAll godoc
// @Summary Listar usuários
// @Description Retorna todos os usuários
// @Tags usuarios
// @Produce json
// @Success 200 {array} models.Usuario
// @Failure 500 {object} map[string]string
// @Router /usuarios [get]
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

// FindByID godoc
// @Summary Buscar usuário por ID
// @Description Retorna um usuário específico pelo seu ID
// @Tags usuarios
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /usuarios/{id} [get]
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

// Update godoc
// @Summary Atualizar usuário
// @Description Atualiza um usuário existente
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param usuario body models.Usuario true "Dados atualizados do usuário"
// @Success 200 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /usuarios/{id} [put]
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

// Delete godoc
// @Summary Excluir usuário
// @Description Exclui um usuário existente
// @Tags usuarios
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /usuarios/{id} [delete]
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
