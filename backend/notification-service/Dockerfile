FROM public.ecr.aws/docker/library/golang:1.26-alpine AS build

WORKDIR /app
COPY backend/proto/gen/go /proto/gen/go
COPY backend/notification-service/go.mod backend/notification-service/go.sum ./
RUN go mod download
COPY backend/notification-service ./
RUN go build -o /bin/notification-service ./cmd/server

FROM public.ecr.aws/docker/library/alpine:3.21

COPY --from=build /bin/notification-service /bin/notification-service
EXPOSE 8085
ENTRYPOINT ["/bin/notification-service"]
