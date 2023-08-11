FROM golang:alpine

RUN apk update && apk add --no-cache git

RUN apk add --no-cache nginx

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go mod tidy

COPY . .

# Build aplikasi
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]