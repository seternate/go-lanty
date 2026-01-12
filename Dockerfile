FROM golang:1.25-alpine AS build

ARG APP_VERSION

RUN apk add --no-cache git

COPY . .

RUN set -e && \
    if [ -z "$APP_VERSION" ]; then \
      VERSION="dev-build-$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"; \
    else \
      VERSION="$APP_VERSION"; \
    fi && \
    CGO_ENABLED=0 go build -o /go/bin/lantyd -ldflags "-X main.AppVersion=$VERSION" ./cmd/lantyd

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/lantyd /lantyd

ENTRYPOINT [ "/lantyd" ]
