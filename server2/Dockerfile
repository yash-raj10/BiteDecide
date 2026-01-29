# FROM golang:1.22.2

# WORKDIR /app

# COPY . .


# RUN go mod tidy
# RUN go build main.go

# EXPOSE 8000

# CMD ["./main"]



FROM golang:1.22.2-alpine AS ash
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go

FROM alpine
WORKDIR /app

COPY --from=ash /app/main .
COPY --from=ash /app/data ./data
EXPOSE 8000

CMD ["./main"]
