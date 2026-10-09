FROM public.ecr.aws/docker/library/golang:1.26-alpine AS build

WORKDIR /app
COPY backend/matching-service/go.mod backend/matching-service/go.sum ./
RUN go mod download
COPY backend/matching-service ./
RUN go build -o /bin/matching-service ./cmd/server

FROM public.ecr.aws/docker/library/alpine:3.21

COPY --from=build /bin/matching-service /bin/matching-service
EXPOSE 8084
ENTRYPOINT ["/bin/matching-service"]
