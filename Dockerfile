FROM golang:1.20 AS build

ARG REST_PORT=8080

RUN apt-get update
RUN apt install -y protobuf-compiler

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY third_party/ ./third_party/

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=direct

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN go mod download
RUN go mod verify
RUN go generate ./third_party/dataway
RUN go build \
        -ldflags "-s -w" \
        -o datatom ./cmd/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/datatom /

EXPOSE ${REST_PORT}

ENTRYPOINT ["/datatom"]