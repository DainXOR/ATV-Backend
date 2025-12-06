FROM golang:1.25.5

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o atv-service .

EXPOSE 8080

ENTRYPOINT [ "./atv-service" ]
