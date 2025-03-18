# Bun step, simply to copy bun binary
FROM oven/bun:1.2.5-debian@sha256:b3af6fbe497a5c451ca283b2097db55fd9d7a2c6afc09acda1ccf10b7553a5ef as bun

# Builder step
FROM golang:1.24-bookworm@sha256:96a23e75187cec3661fc96454cf7adbc5de5f1a2cfea4d5374bed9a464a32287 as builder

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
FROM gcr.io/distroless/static@sha256:95ea148e8e9edd11cc7f639dc11825f38af86a14e5c7361753c741ceadef2167

COPY --from=builder /app/ptht ptht

EXPOSE 8080
ENTRYPOINT ["./ptht", "serve", "--http=0.0.0.0:8080"]
