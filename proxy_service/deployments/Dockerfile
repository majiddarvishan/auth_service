# ---------- STAGE 1: Build ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Set a custom cache directory
# ENV GOCACHE=/go-cache

# Install swag (optional for prod — skip if docs not needed)
# RUN --mount=type=cache,target=/go-cache \
#     env GOCACHE=/go-cache && go install github.com/swaggo/swag/cmd/swag@latest

# RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
# RUN go mod download
RUN --mount=type=cache,target=/go-cache \
    env GOCACHE=/go-cache && go mod download

COPY . .

# Generate Swagger docs (comment this out if not needed)
# RUN swag init

# Install swagger2openapi via npm
# RUN apk add --no-cache nodejs npm && \
# npm install -g swagger2openapi

# Generate Swagger and OpenAPI docs
# RUN swag init -g main.go -o docs && \
# swagger2openapi docs/swagger.json -o docs/openapi3.json


# Build the binary
# RUN go build -o main .
RUN --mount=type=cache,target=/go-cache \
    env GOCACHE=/go-cache && CGO_ENABLED=0 GOOS=linux go build -o main .

# ---------- STAGE 2: Production ----------
FROM alpine:3.21.3

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
# (Optional) Copy docs if using Swagger in prod
COPY --from=builder /app/docs ./docs

EXPOSE 8080
EXPOSE 8443

CMD ["./main"]
