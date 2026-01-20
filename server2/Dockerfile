FROM golang:1.22.2

WORKDIR /app

COPY . .


RUN go mod tidy
RUN go build main.go

EXPOSE 8000

CMD ["./main"]