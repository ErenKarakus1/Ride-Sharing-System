FROM golang:1.26-alpine AS build

WORKDIR /app
COPY backend/proto/gen/go /proto/gen/go
COPY backend/ride-service/go.mod backend/ride-service/go.sum ./
RUN go mod download
COPY backend/ride-service ./
RUN go build -o /bin/ride-service ./cmd/server

FROM alpine:3.21

COPY --from=build /bin/ride-service /bin/ride-service
EXPOSE 8082 9093
ENTRYPOINT ["/bin/ride-service"]
