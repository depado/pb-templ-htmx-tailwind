# Bun step, simply to copy bun binary
FROM oven/bun:1.2.20-debian@sha256:6c0231b1fc1d7eb6820efdf1b0208bdeb097fc4c1c7b9fc5fb6dde7fbf1da6e5 as bun

# Builder step
FROM golang:1.25-bookworm@sha256:74908ad827a5849c557eeca81d46263acf788ead606102d83466f499f83e35b1 as builder

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
FROM gcr.io/distroless/static@sha256:2e114d20aa6371fd271f854aa3d6b2b7d2e70e797bb3ea44fb677afec60db22c

COPY --from=builder /app/ptht ptht

EXPOSE 8080
ENTRYPOINT ["./ptht", "serve", "--http=0.0.0.0:8080"]
