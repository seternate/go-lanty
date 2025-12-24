-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "games" (
    "slug" VARCHAR PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "games";
-- +goose StatementEnd
