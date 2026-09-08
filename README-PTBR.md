# Classificador de Criminalidade[cite: 1]

> 🎯 **Objetivos de Desenvolvimento Sustentável (ODS) Abrangidos:**
> - **ODS 16: Paz, Justiça e Instituições Eficazes** — Promover sociedades pacíficas e inclusivas para o desenvolvimento sustentável, auxiliando na conscientização e redução da criminalidade local por meio do acesso à informação[cite: 1].
> - **ODS 11: Cidades e Comunidades Sustentáveis** — Tornar as cidades e os assentamentos humanos inclusivos, seguros, resilientes e sustentáveis através da geolocalização e transparência de relatos urbanos[cite: 1].

O **Classificador de criminalidade** é uma aplicação de conscientização sobre criminalidade baseada em localização[cite: 1]. Ele exibe relatos de crimes próximos em um mapa interativo, utilizando a localização atual do usuário e um raio de busca configurável[cite: 1]. O backend também disponibiliza endpoints CRUD para relatos de crimes e usuários[cite: 1].

## Funcionalidades principais[cite: 1]

- Mapa interativo alimentado por Leaflet e OpenStreetMap[cite: 1]
- Busca automática de relatos de crimes próximos à localização do usuário[cite: 1]
- Consultas geoespaciais de crimes utilizando o índice `2dsphere` do MongoDB[cite: 1]
- API REST desenvolvida com Go e Gin[cite: 1]
- Documentação da API com Swagger[cite: 1]
- Configuração com Docker Compose para MongoDB, backend e frontend[cite: 1]

## Pré-requisitos[cite: 1]

Para a instalação recomendada, instale:[cite: 1]

- Docker Desktop com Docker Compose[cite: 1]

Para desenvolvimento local sem Docker, instale também:[cite: 1]

- Go 1.26 ou superior[cite: 1]
- Node.js 22 ou superior[cite: 1]
- npm[cite: 1]

## Executar com Docker Compose[cite: 1]

1. Na raiz do repositório, crie o arquivo de ambiente do backend:[cite: 1]

	```powershell
	Copy-Item backend/.env.example backend/.env
	```[cite: 1]

	Atualize os valores em `backend/.env` caso não queira utilizar os padrões de desenvolvimento[cite: 1].

2. Compile e inicie todos os serviços:[cite: 1]

	```powershell
	docker compose up --build
	```[cite: 1]

3. Acesse a aplicação em [http://localhost:3000](http://localhost:3000)[cite: 1].

A API estará disponível em [http://localhost:8080](http://localhost:8080), e sua documentação Swagger em [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)[cite: 1].

Para parar os serviços, pressione `Ctrl+C` ou execute:[cite: 1]

```powershell
docker compose down
```[cite: 1]

Para parar os serviços e remover o volume de dados do MongoDB:[cite: 1]

```powershell
docker compose down -v
```[cite: 1]

## Executar localmente[cite: 1]

O backend local espera que o MongoDB esteja acessível utilizando a string de conexão em `backend/.env`[cite: 1]. Você pode iniciar apenas o MongoDB com o Docker:[cite: 1]

```powershell
Copy-Item backend/.env.example backend/.env
docker compose up mongo
```[cite: 1]

Em um segundo terminal, inicie o backend:[cite: 1]

```powershell
Set-Location backend
MongoDB Shell - [https://www.mongodb.com/try/download/shell](https://www.mongodb.com/try/download/shell)
```[cite: 1]

Em um terceiro terminal, inicie o frontend:[cite: 1]

```powershell
Set-Location frontend
npm install
npm run dev
```[cite: 1]

O servidor de desenvolvimento do Vite normalmente roda em [http://localhost:5173](http://localhost:5173)[cite: 1].

## Rotas da API[cite: 1]

A API está agrupada sob `/api` e atualmente disponibiliza:[cite: 1]

- `GET`, `POST`, `PUT` e `DELETE` `/api/crimes`[cite: 1]
- `GET` `/api/crimes/proximos?lat={latitude}&lng={longitude}&raio={metros}`[cite: 1]
- `GET`, `POST`, `PUT` e `DELETE` `/api/usuarios`[cite: 1]

Utilize o Swagger para visualizar os esquemas completos de requisição e resposta[cite: 1].

## Permissão importante do navegador[cite: 1]

O frontend solicita acesso à API de geolocalização do navegador[cite: 1]. Permita o acesso à localização quando solicitado para que o mapa possa buscar os relatos próximos[cite: 1]. Caso a permissão seja negada, o mapa ainda será carregado, mas os relatos próximos não poderão ser solicitados automaticamente[cite: 1].

## Testes[cite: 1]

Execute os testes do backend com:[cite: 1]

```powershell
Set-Location backend
Postman - [https://www.postman.com/downloads/](https://www.postman.com/downloads/)
```[cite: 1]

## Estrutura do projeto[cite: 1]

```text
backend/     API em Go, serviços, repositórios, modelos e documentação Swagger
frontend/    Aplicação web com Vite e mapa Leaflet
docker-compose.yml
             Serviços do MongoDB, backend e frontend
```[cite: 1]