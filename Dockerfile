FROM golang:1.19-buster AS builder
WORKDIR /go/src/
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /go/bin/app .

FROM gcr.io/distroless/static
COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]
