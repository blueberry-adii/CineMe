BINARY_PATH=build
BINARY_NAME=app
TARGET=$(BINARY_PATH)/$(BINARY_NAME)
PID_FILE=/tmp/$(BINARY_NAME).pid

SRC=cmd/main.go
SOURCES=$(shell find . -name "*.go")

.PHONY: all build run test clean watch stop

all: build

$(TARGET): $(SOURCES)
	@mkdir -p $(BINARY_PATH)
	@echo "Files changed, rebuilding..."
	go build -o $(TARGET) $(SRC)

build: $(TARGET)

# Starts the server and saves its process ID
run: build
	@make stop
	./$(TARGET) & echo $$! > $(PID_FILE)

# Safely kills the server if it's running
stop:
	@if [ -f $(PID_FILE) ]; then \
		kill $$(cat $(PID_FILE)) 2>/dev/null || true; \
		rm $(PID_FILE); \
	fi

clean: stop
	rm -rf $(BINARY_PATH)

watch:
	@echo "Watching for changes (Ctrl+C to stop)..."
	@make run # Initial start
	@while true; do \
		if make -q build; then \
			sleep 1; \
		else \
			echo "Change detected! Restarting..."; \
			make run; \
		fi; \
	done