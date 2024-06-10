FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build
EXPOSE 3333

CMD ["./main"]