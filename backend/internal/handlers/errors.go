package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Códigos padronizados para categorização de erros pela API
const (
	CodigoEntradaInvalida       = "ENTRADA_INVALIDA"
	CodigoRegistroNaoEncontrado = "REGISTRO_NAO_ENCONTRADO"
	CodigoFalhaInterna          = "FALHA_INTERNA"
)

// ErroResposta define o formato JSON consistente retornado em caso de erro.
type ErroResposta struct {
	Codigo   string `json:"codigo"`
	Mensagem string `json:"mensagem"`
	Detalhes string `json:"detalhes,omitempty"`
	Erro     string `json:"erro"` // Campo mantido para retrocompatibilidade
}

// RespostaEntradaInvalida retorna HTTP 400 com código ENTRADA_INVALIDA.
func RespostaEntradaInvalida(c *gin.Context, mensagem string, detalhes ...string) {
	var det string
	if len(detalhes) > 0 {
		det = detalhes[0]
	}
	c.JSON(http.StatusBadRequest, ErroResposta{
		Codigo:   CodigoEntradaInvalida,
		Mensagem: mensagem,
		Detalhes: det,
		Erro:     mensagem,
	})
}

// RespostaNaoEncontrado retorna HTTP 404 com código REGISTRO_NAO_ENCONTRADO.
func RespostaNaoEncontrado(c *gin.Context, mensagem string) {
	c.JSON(http.StatusNotFound, ErroResposta{
		Codigo:   CodigoRegistroNaoEncontrado,
		Mensagem: mensagem,
		Erro:     mensagem,
	})
}

// RespostaFalhaInterna retorna HTTP 500 com código FALHA_INTERNA.
func RespostaFalhaInterna(c *gin.Context, mensagem string, detalhes ...string) {
	var det string
	if len(detalhes) > 0 {
		det = detalhes[0]
	}
	c.JSON(http.StatusInternalServerError, ErroResposta{
		Codigo:   CodigoFalhaInterna,
		Mensagem: mensagem,
		Detalhes: det,
		Erro:     mensagem,
	})
}
