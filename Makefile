.DEFAULT_GOAL := build

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: templ
templ: # Generate templ files
	templ generate

.PHONY: assets
assets: # Generate CSS based on templ files
	bun run tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css

.PHONY: embed
embed: templ assets # Embed generated assets
	go generate ./...

.PHONY: build
build: embed # Generate, embed and build with proper flags
	go build -ldflags="-s -w" -o ptht

.PHONY: run
run: templ assets # Run the server
	go run main.go serve

.PHONY: dev
dev: # Run in dev mode with file watching
	wgo -file=.go -file=.templ -xfile=_templ.go templ generate :: \
		bun run tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css :: \
		go run main.go serve

.PHONY: templdev
templdev: # Run with templ's hot reload proxy
	templ generate --watch --proxy="http://127.0.0.1:8090" --cmd="$(MAKE) --no-print-directory templrun"

.PHONY: .templrun
.templrun: assets # Not invokable because there's no point in running this without templ watch
	go run main.go serve --http=127.0.0.1:8090

# .PHONY: templdevwgo
# templdevwgo:
# 	wgo -exit -file=.go-xfile=_templ.go templ generate --watch --proxy="http://127.0.0.1:8090" --cmd="$(MAKE) --no-print-directory templrun"

