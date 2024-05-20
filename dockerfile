FROM golang:1.22.2

ENV TODO_PORT 7666
ENV TODO_DBFILE ./scheduler.db
ENV TODO_PASSWORD 12345

WORKDIR /app_bin

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY app ./app
COPY tests ./tests
COPY web ./web

EXPOSE 7666
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./my_app

CMD ["./my_app"]