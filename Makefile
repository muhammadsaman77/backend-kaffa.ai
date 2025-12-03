migrate_generate:
	migrate create -ext sql -dir  db/migrations -seq $(NAME)
migrate_up:
	migrate -path db/migrations -database postgres://kaffa_user:kaffa_password@localhost:5436/kaffa_db?sslmode=disable up
migrate_down:
	migrate -path db/migrations -database postgres://kaffa_user:kaffa_password@localhost:5436/kaffa_db?sslmode=disable down
migrate_version:
	migrate -path db/migrations -database postgres://kaffa_user:kaffa_password@localhost:5436/kaffa_db?sslmode=disable version
migrate_force:
	migrate -path db/migrations -database postgres://kaffa_user:kaffa_password@localhost:5436/kaffa_db?sslmode=disable force $(VERSION)
