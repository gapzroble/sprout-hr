FROM golang:1.23.6

WORKDIR /app

COPY . ./

RUN make build

EXPOSE 3000

CMD [ "./go-app" ]
