FROM golang:1.22-alpine3.20 AS builder

ARG TARGETARCH
ARG TARGETOS

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-${TARGETARCH}.tar.gz | tar xvz

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY db/migration ./migration

EXPOSE 80
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
