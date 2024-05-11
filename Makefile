.DEFAULT_GOAL := build

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

templ: # Generate templ files
	templ generate

assets: templ # Generate CSS based on templ files
	bun run tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css

generate: assets # Embed generated assets
	go generate ./...

build: generate # Generate, embed and build with proper flags
	go build -ldflags="-s -w" -o ptht

run: assets # Run the server
	go run main.go serve

dev: # Run in dev mode with file watching
	wgo -file=.go -file=.templ -xfile=_templ.go make run
