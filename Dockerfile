# Bun step, simply to copy bun binary
FROM oven/bun:1.1.8-debian as bun

# Builder step
FROM golang:1.22-bookworm as builder

# Setup bun
COPY --chown=root:root --from=bun /usr/local/bin/bun /root/.bun/bin/
ENV PATH="${PATH}:/root/.bun/bin"

# Change workdir
WORKDIR /app

# Install bun deps
COPY bun.lockb package.json ./
RUN bun install

# Install go deps and templ
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate templ, build CSS and embed assets
COPY . .
RUN templ generate
RUN bun run tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css
RUN go generate ./...

# Build binary
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags="-s -w" -o htmxtodo

# Serve step
FROM gcr.io/distroless/static

COPY --from=builder /app/htmxtodo htmxtodo

EXPOSE 8080
ENTRYPOINT ["./htmxtodo", "serve", "--http=0.0.0.0:8080"]
