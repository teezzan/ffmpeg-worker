
FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /worker_bin cmd/worker/main.go


CMD [ "/worker_bin" ]