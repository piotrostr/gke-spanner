FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/spanner-go-experiment .


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Europe/Warsaw /usr/share/zoneinfo/Europe/Warsaw
ENV TZ Europe/Warsaw

WORKDIR /app
COPY --from=builder /app/spanner-go-experiment /app/spanner-go-experiment

CMD ["./spanner-go-experiment"]
