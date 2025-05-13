FROM golang:1.23.6 AS builder

WORKDIR /app

COPY /cmd/main.go ./cmd/

RUN go mod init go-app

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app ./cmd/main.go

FROM scratch

WORKDIR /

COPY --from=builder /go-app .

EXPOSE 3000

CMD [ "/go-app" ]
