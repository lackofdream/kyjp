FROM golang:alpine AS builder

COPY . .

RUN go build -ldflags "-s -w" -o /a.out .

FROM apline

COPY --from=builder /a.out /a.out

ENTRYPOINT ["/a.out"]