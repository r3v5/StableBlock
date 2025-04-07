# Set default migration path
MIGRATIONS_DIR=database/migrations

# Run all migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

# Roll back last migration
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1

# Roll back all migrations
migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down

# Re-run migrations from scratch
migrate-restart:
	make migrate-reset
	make migrate-up

# Show current migration version
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version

# Create a new migration: make migrate-new name=add_users_table
migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name)
