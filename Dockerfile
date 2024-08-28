FROM golang:1.22-alpine

WORKDIR /server

RUN apk update && apk add --no-cache git && \
    go install github.com/cosmtrek/air@latest


COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV PORT=3000
ENV GIN_MODE=release

RUN go build -o main .

EXPOSE 3000


COPY air.sh /server/air.sh
RUN chmod +x /server/air.sh

CMD ["/server/air.sh"]