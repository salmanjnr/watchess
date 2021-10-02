db-user = dev
db-pass = devpass

build:
	find cmd/web/ -name '*.go' -not -path '*_test.go' | xargs go build -o bin/server

run:
	find cmd/web/ -name '*.go' -not -path '*_test.go' | xargs go run

format:
	find -name '*.go' | xargs gofmt -w -l

migrateup:
	migrate -path db/migration -database "mysql://$(db-user):$(db-pass)@/watchess" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://$(db-user):$(db-pass)@/watchess" -verbose down

dbseed:
	mysql -h localhost -u $(db-user) -p$(db-pass) watchess < db/seed/users.sql
	mysql -h localhost -u $(db-user) -p$(db-pass) watchess < db/seed/tournaments.sql

test:
	go test -v ./...
