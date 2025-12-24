FROM golang:1.25.5 AS installer

RUN git clone --branch v3.26.0 --depth 1 https://github.com/pressly/goose.git && \
    cd goose && \
    go mod tidy && \
    CGO_ENABLED=0 go build -tags='no_clickhouse no_libsql no_sqlite3 no_mssql no_vertica no_mysql no_ydb' -ldflags="-s -w" -o /go/bin/goose ./cmd/goose

FROM scratch

COPY *.sql /
COPY --from=installer /go/bin/goose /goose

ENTRYPOINT [ "/goose", "up" ]
