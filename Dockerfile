FROM golang:1.17.8 as builder
WORKDIR /go/github.com/GeneralKenobi/mailman

COPY go.mod go.sum ./
RUN go mod download

COPY pkg pkg
COPY internal internal
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o mailman cmd/mailman/main.go


FROM alpine:3.15
WORKDIR /opt/mailman
COPY --from=builder /go/github.com/GeneralKenobi/mailman/mailman mailman
RUN chmod 755 mailman

EXPOSE 8080
ENTRYPOINT ["/opt/mailman/mailman"]
