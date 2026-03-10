.PHONY: all build clean run compressed help css

APP_NAME = go-samba4
DATE     = $(shell date +%Y-%m-%d\ %H:%M:%S)
VERSION  = v1.1.2
PREFIX   = go-samba4/internal/buildinfo
LDFLAGS  = -X '$(PREFIX).Version=$(VERSION)' -X '$(PREFIX).BuildDate=$(DATE)'
FLAGS    = -v -ldflags="-s -w $(LDFLAGS)"

all: build

css:
	@echo "Compiling Tailwind CSS..."
	@if [ ! -d node_modules ]; then npm install > /dev/null; fi
	npx --yes @tailwindcss/cli -i web/static/css/input.css -o web/static/css/app.css --minify
	@echo "CSS compiled."

build: clean css
	@echo "Building $(APP_NAME) $(VERSION)..."
	go build $(FLAGS) -o $(APP_NAME) .
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

run:
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME) serve

certs:
	@echo "Generating SSL certificates..."
	mkdir -p ssl
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
		-keyout ssl/server.key -out ssl/server.crt \
		-subj "/C=BR/ST=SP/L=Sao Paulo/O=Development/CN=localhost"

help:
	@echo "Available commands:"
	@echo "  make build      - Build the application (stripped binary)"
	@echo "  make compressed - Build and compress the binary using UPX"
	@echo "  make css        - Compile Tailwind CSS"
	@echo "  make clean      - Remove the generated binary"
	@echo "  make run        - Build and start the 'serve' command"
