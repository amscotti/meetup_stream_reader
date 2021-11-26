FROM golang:1.17 AS builder

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /meetup_stream_reader .


FROM scratch

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /meetup_stream_reader /meetup_stream_reader

ENTRYPOINT ["/meetup_stream_reader"]