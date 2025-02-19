
# this build target is used to ensure all the golang code in all subdirectories compiles
build:
	go build -o bin/jam ./cmd/jam

build-all-subpackages:
	go build ./...

install:
	go install ./cmd/jam

run:
	go run ./cmd/jam run

generate-db:
	go run ./models/tools/dbgen

# Run goimports against code
goimports:
	goimports -w .

# Run go vet against code
vet:
	go vet ./...

# lint the code
lint:
	golangci-lint run
	
# tidy the go modules
tidy: 
	go mod tidy

migrate:
	go run ./cmd/jam migrate


# launches the godoc server on port 6060
godocs:
	@echo "open http://localhost:6060/pkg/bhbdev/jam/"
	godoc -http=:6060

postcss:
	npx postcss ./assets/css/styles.css -o ./assets/css/styles.css --use autoprefixer --use cssnano

buildcss:
	npx @tailwindcss/cli -i ./app/styles.css -o ./assets/css/styles.css
watchcss:
	npx @tailwindcss/cli -i ./app/styles.css -o ./assets/css/styles.css --watch

.PHONY: build build-all-subpackages install run generate-db goimports vet lint tidy migrate godocs