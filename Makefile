BINARY_NAME=main
SRC_DIR=.
TEMPL_CMD=templ generate
WATCH_DIRS=$(SRC_DIR)/*.go $(SRC_DIR)/*.templ $(SRC_DIR)/*.html $(SRC_DIR)/*.css
PID_FILE=server.pid

build:
	go build -o $(BINARY_NAME) $(SRC_DIR)/main.go

run: build
	@if [ -f $(PID_FILE) ]; then \
        kill `cat $(PID_FILE)`; \
        rm $(PID_FILE); \
    fi
	./$(BINARY_NAME) & echo $$! > $(PID_FILE)

generate:
	$(TEMPL_CMD)

watch:
	@echo "Watching for changes..."
	@fswatch -o $(WATCH_DIRS) --exclude $(BINARY_NAME) --exclude $(PID_FILE) | xargs -I{} sh -c 'make generate && make run'

.PHONY: build run generate watch