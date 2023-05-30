FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build cmd/server/server.go -o /ssd-bot-go

#EXPOSE 8080

CMD [ "/ssd-bot-go" ]
