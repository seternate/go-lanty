-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "assets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "url" VARCHAR NOT NULL,
    "size" BIGINT NOT NULL,
    "checksum" VARCHAR NOT NULL,
    "algorithm" VARCHAR NOT NULL,
    "mime_type" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "assets";
-- +goose StatementEnd
