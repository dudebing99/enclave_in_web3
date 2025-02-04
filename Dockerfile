FROM golang:1.19-alpine AS builder

# Install dependencies for building (if needed, e.g., git for fetching dependencies)
RUN apk add --no-cache gcc libc-dev make

WORKDIR /enclave_in_web3

COPY . .

RUN go mod tidy

RUN go build -o bin/enclave-keeper ./cmd/keeper/main.go

FROM alpine:latest

WORKDIR /enclave_in_web3

# Install necessary dependencies for running the app (e.g., libc)
#RUN apk add --no-cache libc6-compat

COPY --from=builder /enclave_in_web3/bin/enclave-keeper ./
COPY --from=builder /enclave_in_web3/conf/application_keeper.yml ./conf/application.yml

EXPOSE 10000

#CMD ["sleep", "infinity"]
CMD ["/enclave_in_web3/enclave-keeper", "-config", "/enclave_in_web3/conf/application.yml"]
