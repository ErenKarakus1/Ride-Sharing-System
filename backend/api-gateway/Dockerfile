FROM mirror.gcr.io/library/golang:1.26-alpine AS build

WORKDIR /app
COPY backend/proto/gen/go /proto/gen/go
COPY backend/api-gateway/go.mod backend/api-gateway/go.sum ./
RUN go mod download
COPY backend/api-gateway ./
RUN go build -o /bin/api-gateway ./cmd/server

FROM mirror.gcr.io/library/alpine:3.21

COPY --from=build /bin/api-gateway /bin/api-gateway
EXPOSE 8088
ENTRYPOINT ["/bin/api-gateway"]
