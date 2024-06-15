FROM golang:latest

RUN mkdir -p /home/app/build

WORKDIR /home/app

COPY . .

RUN go mod tidy

RUN go build -o ./build/app ./cmd/app/*.go

EXPOSE 3000

CMD ["./build/app"]
