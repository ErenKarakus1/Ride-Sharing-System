FROM public.ecr.aws/docker/library/golang:1.26-alpine AS build

WORKDIR /app
COPY backend/proto/gen/go /proto/gen/go
COPY backend/user-service/go.mod backend/user-service/go.sum ./
RUN go mod download
COPY backend/user-service ./
RUN go build -o /bin/user-service ./cmd/server

FROM public.ecr.aws/docker/library/alpine:3.21

COPY --from=build /bin/user-service /bin/user-service
EXPOSE 8080 9091
ENTRYPOINT ["/bin/user-service"]
