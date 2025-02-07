# pinger/Dockerfile
FROM golang:1.21-alpine  # Версия 1.23 не существует, лучше 1.21

WORKDIR /app

COPY go.mod ./
RUN go mod tidy  

COPY . .

RUN go build -o pinger .

CMD ["./pinger"]

