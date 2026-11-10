# Start from the latest golang base image
FROM golang:latest AS builder
ARG CURRENT_DIRECTORY

# Set the Current Working Directory inside the container
WORKDIR /app

# Set up Healthcheck 
COPY --from=ghcr.io/grpc-ecosystem/grpc-health-probe:v0.4.38 /ko-app/grpc-health-probe /bin/grpc_health_probe

# Copy go mod and sum files
COPY proto/go.mod proto/go.sum ./proto/
COPY ${CURRENT_DIRECTORY}/go.mod ${CURRENT_DIRECTORY}/go.sum ./${CURRENT_DIRECTORY}/

# Initialize work to combine proto and current directory dependencies
RUN go work init ./proto ./${CURRENT_DIRECTORY}

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY proto/ ./proto/
COPY ${CURRENT_DIRECTORY}/ ./${CURRENT_DIRECTORY}/

# Build the Go app
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o main ./${CURRENT_DIRECTORY}


######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

# Healthcheck
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD /bin/grpc_health_probe -addr=localhost:8000 || exit 1

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]