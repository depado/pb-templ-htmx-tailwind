# Bun step, simply to copy bun binary
FROM oven/bun:1.3.10-debian@sha256:367842b35abbdf23f39e23c71f3a08eee940ff2679a14e08a5afcf4a1436cd89 as bun

# Builder step
FROM golang:1.26-bookworm@sha256:2a0ba12e116687098780d3ce700f9ce3cb340783779646aafbabed748fa6677c as builder

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
FROM gcr.io/distroless/static@sha256:eca24e67792afe660769e84c85332d9939772e7f76071597d179af96ac4e9e4f

COPY --from=builder /app/ptht ptht

EXPOSE 8080
ENTRYPOINT ["./ptht", "serve", "--http=0.0.0.0:8080"]
