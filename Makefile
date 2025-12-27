migrate_create:
	migrate create -seq  -ext sql -dir db/migrations

migrate_up:
	migrate -path db/migrations -database $(DATABASE_URL) up


generate:
	sqlc -src ./db/sqlc -dst ./db/sqlc -config ./db/sqlc/sqlc.yaml

.PHONY: migrate_create migrate_up