FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /ssd-bot-go cmd/server/server.go

#EXPOSE 8080

CMD [ "/ssd-bot-go" ]
