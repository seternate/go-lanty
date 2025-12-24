-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "game_assets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "game_slug" VARCHAR NOT NULL REFERENCES "games"("slug") ON DELETE CASCADE,
    "asset_id" UUID NOT NULL REFERENCES "assets"("id") ON DELETE CASCADE,
    "role" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("game_slug", "role")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "game_assets";
-- +goose StatementEnd
