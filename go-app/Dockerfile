FROM golang:1.23.6 AS builder

WORKDIR /app

COPY . ./

RUN make build

FROM scratch

WORKDIR /

COPY --from=builder /app/go-app .

EXPOSE 3000

CMD [ "/go-app" ]
