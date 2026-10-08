FROM golang:1.26-alpine AS build

WORKDIR /app
COPY backend/location-service/go.mod backend/location-service/go.sum ./
RUN go mod download
COPY backend/location-service ./
RUN go build -o /bin/location-service ./cmd/server

FROM alpine:3.21

COPY --from=build /bin/location-service /bin/location-service
EXPOSE 8083
ENTRYPOINT ["/bin/location-service"]
