FROM goLang:1.22.0 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /api /api

EXPOSE 3500

ENTRYPOINT [ "./api" ]