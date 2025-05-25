FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN echo "" > .env
RUN make prepare

COPY . .

RUN make build

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/migrations /migrations
COPY --from=build-stage /app/gofemart-service /gofemart-service
COPY --from=build-stage /app/migrator-service /migrator-service

USER nonroot:nonroot

ENTRYPOINT ["/gofemart-service"]