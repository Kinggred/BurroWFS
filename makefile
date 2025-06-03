MIGRATE_BIN=cmd/migrate_up.go
.PHONY: test, migrate, build, clean

test:
	go test -cover ./...

build:
	go build -o migrate-tool $(MIGRATE_BIN)

migrate: build
	go run $(MIGRATE_BIN)
	rm -f migrate-tool
