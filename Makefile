build:
#	./tailwindcss -i views/css/styles.css -o public/styles.css
	@templ generate
	@go build -o main main.go 

test:
	@go test -v ./...
	
run: build
	@./main

# tailwind:
# 	@./tailwindcss -i views/css/styles.css -o public/styles.css --watch

templ:
	@templ generate -watch -proxy=http://localhost:8008

# migration:
# 	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

# migrate-up:
# 	@go run cmd/migrate/main.go up

# migrate-down:
# 	@go run cmd/migrate/main.go down

watch:
	@air &
	@templ &
	@tailwind &