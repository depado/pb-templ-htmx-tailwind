# Bun step, simply to copy bun binary
FROM oven/bun:1.2.13-debian@sha256:0538136de28963e645f42a9d99da3091c20144a058f27af8fc0dd996ccb4a745 as bun

# Builder step
FROM golang:1.24-bookworm@sha256:89a04cc2e2fbafef82d4a45523d4d4ae4ecaf11a197689036df35fef3bde444a as builder

# Setup bun
COPY --chown=root:root --from=bun /usr/local/bin/bun /root/.bun/bin/
ENV PATH="${PATH}:/root/.bun/bin"

# Change workdir
WORKDIR /app

# Install bun deps
COPY bun.lockb package.json ./
RUN bun install

# Install go deps and templ
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache \
    go mod download
RUN go mod verify
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go install github.com/a-h/templ/cmd/templ@latest

# Generate templ, build CSS and embed assets
COPY . .
RUN templ generate
RUN bun run tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go generate ./...

# Build binary
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go build -ldflags="-s -w" -o ptht

# Serve step
FROM gcr.io/distroless/static@sha256:d9f9472a8f4541368192d714a995eb1a99bab1f7071fc8bde261d7eda3b667d8

COPY --from=builder /app/ptht ptht

EXPOSE 8080
ENTRYPOINT ["./ptht", "serve", "--http=0.0.0.0:8080"]
