FROM golang:1.21-alpine AS builder


ENV GOPATH=/

COPY ./ ./
RUN go mod download && go mod verify
RUN go build -o banner-service ./cmd/main.go

CMD [ "./banner-service" ]
