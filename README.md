# AEP-1B-2026-2

## Classificador de criminalidade

Classificador de criminalidade é uma aplicação de conscientização sobre crimes baseada em localização. Ela exibe ocorrências de crimes próximas em um mapa interativo, utilizando a localização atual do usuário e um raio de busca configurável. O backend também fornece endpoints CRUD para ocorrências de crimes e usuários.

## Objetivos de Desenvolvimento Sustentável

- **ODS 11 — Cidades e Comunidades Sustentáveis**
- **ODS 16 — Paz, Justiça e Instituições Eficazes**

## Funcionalidades

- Mapa interativo com Leaflet e tiles do OpenStreetMap.
- Busca de ocorrências próximas à localização do usuário, com raio configurável.
- Camada de calor para visualizar a concentração das ocorrências.
- API REST em Go e Gin, com persistência no MongoDB.
- Documentação interativa da API com Swagger.
- Ambiente de desenvolvimento com Docker Compose.

## Requisitos

Para executar a aplicação completa, instale Docker Desktop com Docker Compose.

Para desenvolvimento local, também são necessários Go 1.26.6 ou superior, Node.js 22 ou superior e npm. O MongoDB pode ser executado pelo Docker Compose.

## Executar com Docker Compose

Na raiz do repositório, crie a configuração local do backend:

```powershell
Copy-Item backend/.env.example backend/.env
```

O arquivo de exemplo configura as credenciais de desenvolvimento e a URI usada pelo backend no Compose. Altere-as se necessário. Em seguida, inicie os serviços:

```powershell
docker compose up --build
```

Acesse:

- Frontend: [http://localhost:3000](http://localhost:3000)
- API: [http://localhost:8080/api](http://localhost:8080/api)
- Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

O backend executa as migrations pendentes ao iniciar. Para executá-las sem iniciar a API:

```powershell
docker compose build backend
docker compose run --rm backend ./server -migrate-only
```

Para reverter a última migration:

```powershell
docker compose run --rm backend ./server -migrate-down
```

Pare os serviços com `Ctrl+C` ou, em outro terminal, execute `docker compose down`. Para também apagar os dados persistidos do MongoDB, use `docker compose down -v`.

## Desenvolvimento local

Crie o arquivo de ambiente e inicie o MongoDB:

```powershell
Copy-Item backend/.env.example backend/.env
docker compose up -d mongo
```

O `.env.example` usa o hostname `mongo`, válido entre contêineres. Para executar o backend diretamente no Windows, configure a URI para apontar para a porta publicada em localhost:

```powershell
Set-Location backend
$env:MONGODB_URI = "mongodb://aep_root:change_me@localhost:27017/AEP?authSource=admin"
go run .
```

O servidor da API inicia em [http://localhost:8080](http://localhost:8080). Para executar o frontend, use outro terminal:

```powershell
Set-Location frontend
npm install
npm run dev
```

O Vite disponibiliza a aplicação normalmente em [http://localhost:5173](http://localhost:5173). A configuração do Vite encaminha as chamadas `/api` para o backend local.

Para interromper o MongoDB, execute `docker compose down` na raiz do projeto.

## Testes

Com o MongoDB do Compose iniciado, execute os testes do backend:

```powershell
Set-Location backend
go test ./...
```

Os testes de unidades não precisam do MongoDB. Os testes de integração tentam conectar ao MongoDB local usando as credenciais de desenvolvimento e são ignorados quando o banco está indisponível. Para usar outra instância, defina `MONGODB_URI` antes de executar os testes.

## API

Todas as rotas da API usam o prefixo `/api`.

| Recurso | Rotas |
| --- | --- |
| Crimes | `GET/POST /api/crimes`, `GET/PUT/DELETE /api/crimes/{id}` |
| Busca por proximidade | `GET /api/crimes/proximos` |
| Usuários | `GET/POST /api/usuarios`, `GET/PUT/DELETE /api/usuarios/{id}` |

A busca por proximidade aceita `lat` e `lng` (obrigatórios), `raio` (metros, padrão 1000), e filtros opcionais `tipo`, `data_inicio` e `data_fim`. As datas devem estar no formato RFC3339, por exemplo `2026-08-28T12:00:00Z`. O raio informado na interface está em quilômetros e é convertido para metros antes da chamada à API.

Exemplo:

```text
GET /api/crimes/proximos?lat=-23.55&lng=-46.63&raio=500&tipo=roubo&data_inicio=2026-08-28T12:00:00Z
```

Consulte o Swagger para os formatos completos de requisição e resposta.

## Localização no navegador

O mapa pede permissão para acessar a localização do navegador e fazer a busca inicial. Se a permissão for negada ou a localização não estiver disponível, ainda é possível abrir o mapa; use **Usar minha localização** ou permita o acesso para buscar ocorrências próximas.

## Estrutura

```text
backend/       API Go, modelos, handlers, serviços, repositórios, migrations e Swagger
frontend/      Aplicação web com Vite, Leaflet e camada de calor
docker-compose.yml
               Configuração do MongoDB, backend e frontend
```
