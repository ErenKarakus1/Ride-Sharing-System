FROM mirror.gcr.io/library/golang:1.26-alpine AS build

WORKDIR /app
COPY backend/proto/gen/go /proto/gen/go
COPY backend/payment-service/go.mod backend/payment-service/go.sum ./
RUN go mod download
COPY backend/payment-service ./
RUN go build -o /bin/payment-service ./cmd/server

FROM mirror.gcr.io/library/alpine:3.21

COPY --from=build /bin/payment-service /bin/payment-service
EXPOSE 8087
ENTRYPOINT ["/bin/payment-service"]
