APP_NAME := my-app
DB_URL := postgres://postgres:postgres@localhost:5432/peliculas_db?sslmode=disable

.PHONY: all generate migrate apply status build test clean

all: build

generate:
	@sqlc generate
	# @templ generate

# Genera una migración: make migrate name=nombre_de_la_migracion
migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url "docker://postgres/15/dev?search_path=public"

# Aplica las migraciones pendientes
apply:
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

# Muestra el estado de las migraciones
status:
	atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

# Construye el binario de la aplicación
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# Ejecución de tests automatizando Docker y las migraciones de las filminas
test: generate build
	@echo "-> Limpiando entorno..."
	@docker compose down -v
	@echo "-> Levantando PostgreSQL..."
	@docker compose up -d
	@echo "-> Esperando a la base de datos..."
	@sleep 3
	@echo "-> Aplicando migraciones pendientes..."
	@atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"
	@echo "-> Corriendo pruebas..."
	@go test -v ./...
	@echo "-> Limpiando entorno..."
	@docker compose down -v

test-api: generate build
	@echo "-> Limpiando entorno previo..."
	@docker compose down -v
	@echo "-> Levantando PostgreSQL en Docker..."
	@docker compose up -d
	@echo "-> Esperando a que PostgreSQL esté listo..."
	@sleep 3
	@echo "-> Aplicando migraciones con Atlas..."
	@atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"
	@echo "-> Iniciando el servidor Go en segundo plano..."
	@./tmp/$(APP_NAME) & echo $$! > tmp/server.pid
	@sleep 2
	@echo "-> Ejecutando pruebas de integración con Hurl..."
	@hurl --test requests.hurl; EXIT_CODE=$$?; \
	echo "-> Apagando servidor Go..."; \
	kill $$(cat tmp/server.pid) && rm tmp/server.pid; \
	echo "-> Limpiando contenedores Docker..."; \
	docker compose down -v; \
	exit $$EXIT_CODE

clean:
	@rm -rf tmp