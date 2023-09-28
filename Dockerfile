FROM golang:1.21.1 as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/anti-spam cmd/anti-spam/main.go

FROM gcr.io/distroless/base-debian11

COPY --from=builder /build/bin/anti-spam /anti-spam

ENTRYPOINT ["/anti-spam"]