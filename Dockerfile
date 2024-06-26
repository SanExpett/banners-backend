FROM golang:1.21.1-alpine3.18 as build

WORKDIR /var/backend

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download
RUN go build -o main ./cmd/app/main.go

#=========================================================================================
FROM alpine:3.18 as production

WORKDIR /var/backend
COPY --from=build /var/backend/main main
COPY --from=build /go/bin/migrate migrate

RUN mkdir -p /var/log/backend
COPY db/migrations db/migrations

ENV ALLOW_ORIGIN=localhost:3000
ENV PORT_BACKEND=8080
ENV PATH_TO_ROOT=/var/backend
ENV OUTPUT_LOG_PATH=/var/log/backend/logs.json
ENV ERROR_OUTPUT_LOG_PATH=/var/log/backend/err_logs.json

ENV URL_DATA_BASE=postgres://postgres:postgres@localhost/banners?sslmode=disable
ENV SCHEMA=http://

EXPOSE 8080

ENTRYPOINT ./main