# --- START

FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod verify

RUN go build -o main .

FROM ubuntu:22.04 AS runner

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 3000

RUN ls -l /app

CMD ["/app/main"]

# --- EOF
