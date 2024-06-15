FROM golang:alpine

WORKDIR /app

RUN apk upgrade && apk add git

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air", "-c", ".air.toml"]
