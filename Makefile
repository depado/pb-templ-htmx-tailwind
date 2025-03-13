.DEFAULT_GOAL := build

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: # Installs templ and and bun dependencies
	go install github.com/a-h/templ/cmd/templ@latest
	bun install

.PHONY: templ
templ: # Generate templ files
	templ generate

.PHONY: assets
assets: # Generate CSS based on templ files
	bun run tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css

.PHONY: build
build: templ assets # Generate, embed and build with proper flags
	go build -ldflags="-s -w" -o ptht

.PHONY: run
run: templ assets # Run the server
	go run main.go serve

.PHONY: lint
lint:
	golangci-lint run ./...
	deadcode -generated .

# =============
# = Live Mode =
# =============

.PHONY: live/assets
live/assets:
	bun run tailwindcss -w -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css 

.PHONY: live/templ
live/templ:
	templ generate --watch --proxy="http://127.0.0.1:8090" --open-browser=false

.PHONY: live/server
live/server:
	go run github.com/air-verse/air@v1.60.0 \
		--build.bin "go run main.go serve" \
		--build.delay "100" \
		--build.include_ext "go" \
		--build.exclude_dir "node_modules,pb_data,dist" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true \
		--log.main_only true

.PHONY: live/sync
live/sync:
	go run github.com/bokwoon95/wgo@v0.5.6 -file .css templ generate --notify-proxy

.PHONY: live
live: 
	$(MAKE) --no-print-directory -j4 live/assets live/templ live/server live/sync
