FROM golang:alpine AS builder

WORKDIR /app
COPY . .

RUN go build -ldflags "-s -w" -o a.out .

FROM alpine

COPY --from=builder /app/a.out /a.out

ENTRYPOINT ["/a.out"]