FROM golang:alpine AS builder

WORKDIR /avito_hr

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp ./cmd/avito_hr

FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /avito_hr .

EXPOSE 8080

CMD ["./myapp"]