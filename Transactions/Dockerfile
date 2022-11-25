FROM golang

EXPOSE 8080

WORKDIR ./app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go/ ./
COPY streamConfig.json ./

RUN go build -o app
CMD ["./app"]
