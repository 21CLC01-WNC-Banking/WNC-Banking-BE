swagger:
	swag init --parseDependency --parseInternal
wire:
	wire ./internal

MIGRATIONS_DIR=migrations
DATETIME := $(shell date +%Y%m%d%H%M%S)
migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: You must specify a migration name. Usage: make migration name=your_migration_name"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATIONS_DIR)
	@touch $(MIGRATIONS_DIR)/$(DATETIME)_$(name).up.sql
	@touch $(MIGRATIONS_DIR)/$(DATETIME)_$(name).down.sql
	@echo "Created migration files: $(DATETIME)_$(name).up.sql and $(DATETIME)_$(name).down.sql"

# Command to apply migrations (up)
migrate-up:
	@echo "Running migrations (up)..."
	@go run main.go migrate-up

build:
	docker build -t vukhoa23/banking-be:$(tag) .

# Push the Docker image to the registry
push:
	docker push vukhoa23/banking-be:$(tag)
	docker rmi vukhoa23/banking-be:$(tag)

# Build and push the Docker image
build-push: build push