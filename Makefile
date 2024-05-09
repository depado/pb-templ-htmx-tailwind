.DEFAULT_GOAL := build

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

embed: # Embed generated files if any
	go generate ./...

generate: # Generate templ and tailwind css
	templ generate
	bunx tailwindcss -i ./assets/tailwind.css -o ./assets/dist/styles.css

build: generate embed # Generate, embed and build with proper flags
	go build -ldflags="-s -w"

run: generate # Generate templ and tailwind, then run the server
	go run main.go serve

dev: # Watch file changes
	wgo -file=.go -file=.templ -xfile=_templ.go make run
