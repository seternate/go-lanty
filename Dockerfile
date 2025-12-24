FROM golang:1.25-alpine AS build

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/lantyd ./cmd/lantyd

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/lantyd /lantyd

ENTRYPOINT [ "/lantyd" ]
