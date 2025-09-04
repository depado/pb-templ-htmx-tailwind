# Bun step, simply to copy bun binary
FROM oven/bun:1.2.21-debian@sha256:ae21f34625f91285d9098ca52c130845542b2edace3063239c55a8a6d4c1a2df as bun

# Builder step
FROM golang:1.25-bookworm@sha256:6ad9415c04f1cdb7888141cb247cabb28fa74fcc597034494c65d5bc783246a0 as builder

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
FROM gcr.io/distroless/static@sha256:f2ff10a709b0fd153997059b698ada702e4870745b6077eff03a5f4850ca91b6

COPY --from=builder /app/ptht ptht

EXPOSE 8080
ENTRYPOINT ["./ptht", "serve", "--http=0.0.0.0:8080"]
