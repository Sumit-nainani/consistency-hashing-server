# Here Google's distroless image is used to generate the static golang binary.
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

COPY metric ./metric

COPY router ./router

COPY handler ./handler

# Building the Go binary (static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/app /

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app"]
