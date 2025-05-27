# Build stage
FROM golang:1.24-alpine AS builder

# Copy pre-built dependencies
COPY --from=event-processing-service-deps /deps.tar.gz /deps.tar.gz
RUN tar -xzf /deps.tar.gz -C / && rm /deps.tar.gz

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o myapp ./cmd/server

# Runtime stage
FROM scratch
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/config.yaml /config.yaml
EXPOSE 8080
ENTRYPOINT ["/myapp"]