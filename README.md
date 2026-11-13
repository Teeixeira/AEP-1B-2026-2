AEP-1B-2026-2

## ODS

- **ODS 11 - Cidades e Comunidades Sustentáveis**
- **ODS 16 - Paz, Justiça e Instituições Eficazes**

# Classificador de criminalidade

Classificador de criminalidade é uma aplicação de conscientização sobre crimes baseada em localização. Ela exibe ocorrências de crimes próximas em um mapa interativo, utilizando a localização atual do usuário e um raio de busca configurável. O backend também fornece endpoints CRUD para ocorrências de crimes e usuários.

## Main features

- Mapa interativo desenvolvido com Leaflet e OpenStreetMap
- Busca automática por ocorrências de crimes próximas à localização do usuário
- Consultas geoespaciais de crimes utilizando o índice `2dsphere` do MongoDB
- API REST desenvolvida com Go e Gin
- Documentação da API com Swagger
- Configuração com Docker Compose para MongoDB, backend e frontend

## Requirements

Para a configuração recomendada, instale:

- Docker Desktop com Docker Compose

Para desenvolvimento local sem Docker, instale também:

- Go 1.26 ou mais recente
- Node.js 22 ou mais recente
- npm

## Executar com Docker Compose

1. A partir da raiz do repositório, crie o arquivo de ambiente do backend:

	```powershell
	Copy-Item backend/.env.example backend/.env
	```

	Atualize os valores em `backend/.env` caso não queira utilizar os valores padrão de desenvolvimento.

2. Compile e inicie todos os serviços:

	```powershell
	docker compose up --build
	```

3. Acesse a aplicação em [http://localhost:3000](http://localhost:3000).

A API está disponível em [http://localhost:8080](http://localhost:8080), e sua documentação Swagger está disponível em [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

### Executar migrations pelo Docker

Na raiz do repositório, com `backend/.env` configurado, compile a imagem do backend e execute as migrations pendentes:

```powershell
docker compose build backend
docker compose run --rm backend ./server -migrate-only
```

O comando `docker compose run` utiliza as variáveis de `backend/.env` definidas no serviço `backend` e inicia o serviço `mongo` como dependência, se necessário. Para reverter a última migration:

```powershell
docker compose run --rm backend ./server -migrate-down
```

Esses comandos executam o backend em um contêiner temporário e não iniciam a API.

Para interromper os serviços, pressione `Ctrl+C` ou execute:

```powershell
docker compose down
```

Para interromper os serviços e remover o volume de dados do MongoDB:

```powershell
docker compose down -v
```

## Executar localmente

O backend local espera que o MongoDB esteja acessível utilizando a string de conexão presente em `backend/.env`. Você pode iniciar somente o MongoDB com Docker:

```powershell
Copy-Item backend/.env.example backend/.env
docker compose up mongo
```

Em um segundo terminal, inicie o backend:

```powershell
Set-Location backend
MongoDB Shell - https://www.mongodb.com/try/download/shell
```

Em um terceiro terminal, inicie o frontend:

```powershell
Set-Location frontend
npm install
npm run dev
```

O servidor de desenvolvimento do Vite normalmente é executado em [http://localhost:5173](http://localhost:5173).

## Rotas da API

A API está agrupada em `/api` e atualmente fornece:

- `GET`, `POST`, `PUT` e `DELETE` `/api/crimes`
- `GET` `/api/crimes/proximos?lat={latitude}&lng={longitude}&raio={meters}`
- `GET`, `POST`, `PUT` e `DELETE` `/api/usuarios`

Utilize o Swagger para consultar os esquemas completos de requisição e resposta.

## Permissão importante do navegador

O frontend solicita acesso à API de geolocalização do navegador. Permita o acesso à localização quando solicitado para que o mapa possa buscar ocorrências próximas. Caso a permissão seja negada, o mapa ainda será carregado, mas as ocorrências próximas não poderão ser solicitadas automaticamente.

## Testes

Execute os testes do backend com:

```powershell
Set-Location backend
Postman - https://www.postman.com/downloads/
```

## Estrutura do projeto

```text
backend/     API Go, serviços, repositórios, modelos e documentação Swagger
frontend/    Aplicação de mapa com Vite e Leaflet
docker-compose.yml
			 Serviços do MongoDB, backend e frontend
```
