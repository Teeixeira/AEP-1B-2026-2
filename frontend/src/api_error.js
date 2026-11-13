/**
 * Representa um erro padronizado retornado pela API ou ocorrido na requisição.
 */
export class APIError extends Error {
    constructor(codigo, mensagem, status = 0, detalhes = null) {
        super(mensagem);
        this.name = 'APIError';
        this.codigo = codigo;
        this.status = status;
        this.detalhes = detalhes;
    }

    /**
     * Retorna se o erro é elegível para retry automático ou pelo usuário (ex: falhas de rede ou servidor).
     */
    isRetryable() {
        return this.codigo === 'FALHA_INTERNA' || this.codigo === 'ERRO_CONEXAO' || this.status >= 500;
    }

    /**
     * Converte o erro para uma mensagem amigável e consistente para exibição no frontend.
     */
    toUserMessage() {
        switch (this.codigo) {
            case 'ENTRADA_INVALIDA':
                return `Entrada inválida: ${this.message}`;
            case 'REGISTRO_NAO_ENCONTRADO':
                return `Registro não encontrado: ${this.message}`;
            case 'FALHA_INTERNA':
                return `Falha no servidor: ${this.message}. Tente novamente mais tarde.`;
            case 'ERRO_CONEXAO':
                return 'Não foi possível conectar ao servidor. Verifique sua conexão e tente novamente.';
            default:
                return this.message;
        }
    }
}

/**
 * Faz o parse da resposta HTTP de erro, extraindo código, mensagem e detalhes estruturados.
 *
 * @param {Response} response
 * @param {string} fallbackMessage
 * @returns {Promise<APIError>}
 */
export async function parseAPIError(response, fallbackMessage = 'Ocorreu um erro na requisição.') {
    try {
        const body = await response.json();
        const codigo = body.codigo || inferErrorCode(response.status);
        const mensagem = body.mensagem || body.erro || fallbackMessage;
        const detalhes = body.detalhes || null;
        return new APIError(codigo, mensagem, response.status, detalhes);
    } catch {
        const codigo = inferErrorCode(response.status);
        return new APIError(codigo, fallbackMessage, response.status, null);
    }
}

function inferErrorCode(status) {
    if (status === 404) return 'REGISTRO_NAO_ENCONTRADO';
    if (status >= 400 && status < 500) return 'ENTRADA_INVALIDA';
    if (status >= 500) return 'FALHA_INTERNA';
    return 'ERRO_DESCONHECIDO';
}
