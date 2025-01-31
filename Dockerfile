FROM golang:1.19

WORKDIR /enclave_in_web3

COPY . .

RUN go mod tidy

RUN go build -o bin/enclave-client ./cmd/client/main.go

RUN go build -o bin/enclave-keeper ./cmd/keeper/main.go

#CMD ["./myapp"]
