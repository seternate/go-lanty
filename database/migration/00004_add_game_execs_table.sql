-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "game_execs" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "game_slug" VARCHAR NOT NULL REFERENCES "games"("slug") ON DELETE CASCADE,
    "role" VARCHAR NOT NULL,
    "path" VARCHAR NOT NULL,
    "requires_admin" BOOLEAN,
    "format" VARCHAR,
    "arg_seperator" VARCHAR,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("game_slug", "role")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "game_execs";
-- +goose StatementEnd
