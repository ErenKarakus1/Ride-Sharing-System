FROM public.ecr.aws/docker/library/golang:1.26-alpine AS build

WORKDIR /app
COPY backend/proto/gen/go /proto/gen/go
COPY backend/auth-service/go.mod backend/auth-service/go.sum ./
RUN go mod download
COPY backend/auth-service ./
RUN go build -o /bin/auth-service ./cmd/server

FROM public.ecr.aws/docker/library/alpine:3.21

COPY --from=build /bin/auth-service /bin/auth-service
EXPOSE 8081 9092
ENTRYPOINT ["/bin/auth-service"]
