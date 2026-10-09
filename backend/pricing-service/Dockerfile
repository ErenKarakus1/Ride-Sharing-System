FROM mirror.gcr.io/library/golang:1.26-alpine AS build

WORKDIR /app
COPY backend/pricing-service/go.mod backend/pricing-service/go.sum ./
RUN go mod download
COPY backend/pricing-service ./
RUN go build -o /bin/pricing-service ./cmd/server

FROM mirror.gcr.io/library/alpine:3.21

COPY --from=build /bin/pricing-service /bin/pricing-service
EXPOSE 8086
ENTRYPOINT ["/bin/pricing-service"]
