FROM golang:1.22-alpine

WORKDIR /server


COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . .

ENV PORT=3000
ENV GIN_MODE=release

RUN go build -o main .

EXPOSE 3000


CMD ["sh", "-c", "air init && air"]