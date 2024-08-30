FROM golang:1.22.5 AS build

WORKDIR /app

# Set Go proxy to speed up module downloads
ENV GOPROXY=https://proxy.golang.org

# Copy Go module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the source code and build the application with static linking
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /myapp cmd/main.go

# Minimal runtime image
FROM alpine:latest AS run

# Copy the application executable from the build stage
COPY --from=build /myapp /myapp

WORKDIR /app
EXPOSE 3500
CMD ["/myapp"]
