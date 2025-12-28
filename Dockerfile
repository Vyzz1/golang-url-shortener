# Build stage
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/main .


COPY wait-for-it.sh ./wait-for-it.sh

RUN apk add --no-cache dos2unix bash \
    && dos2unix /app/wait-for-it.sh \
    && chmod +x /app/wait-for-it.sh \
    && apk del dos2unix

    
CMD [ "/app/main" ]