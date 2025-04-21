FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /cfmid

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o cfm-server wrapper.go

FROM wishartlab/cfmid:latest
WORKDIR /cfmid

COPY --from=builder /cfmid/cfm-server .

RUN chmod +x cfm-server

EXPOSE 5001
ENTRYPOINT ["./cfm-server"]
