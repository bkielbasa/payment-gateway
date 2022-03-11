FROM golang:1.17-alpine AS builder

RUN mkdir /build/
WORKDIR /build/
COPY go.* ./
RUN go mod download
COPY . ./
WORKDIR /build/
RUN go build -o webapp

FROM alpine
WORKDIR /app
COPY --from=builder /build/webapp /
ENTRYPOINT /webapp