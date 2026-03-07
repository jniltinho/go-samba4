.PHONY: all build clean run compressed help css

APP_NAME = go-samba4

all: build

css:
	@echo "Compiling Tailwind CSS..."
	@if [ ! -d node_modules ]; then npm install > /dev/null; fi
	npx --yes @tailwindcss/cli -i web/static/css/input.css -o web/static/css/app.css --minify
	@echo "CSS compiled."

build: clean css
	@echo "Building $(APP_NAME)..."
	go build -ldflags="-s -w" -o $(APP_NAME) .
	upx --best --lzma $(APP_NAME)
	@echo "Build complete."

compressed: build
	@echo "Compressing $(APP_NAME) with UPX..."
	upx --best --lzma $(APP_NAME)
	@echo "Compression complete."

clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	@echo "Clean complete."

run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME) serve

help:
	@echo "Available commands:"
	@echo "  make build      - Build the application (stripped binary)"
	@echo "  make compressed - Build and compress the binary using UPX"
	@echo "  make css        - Compile Tailwind CSS"
	@echo "  make clean      - Remove the generated binary"
	@echo "  make run        - Build and start the 'serve' command"
