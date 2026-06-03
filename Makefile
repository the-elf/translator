migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up