-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "game_args" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "game_exec_id" UUID NOT NULL REFERENCES "game_execs"("id") ON DELETE CASCADE,
    "role" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "required" BOOLEAN,
    "enabled" BOOLEAN,
    "format" VARCHAR,
    "arg_separator" VARCHAR,
    "arg" VARCHAR NOT NULL,
    "description" VARCHAR,
    "default_string" VARCHAR,
    "default_bool" BOOLEAN,
    "default_int" BIGINT,
    "default_float" DOUBLE PRECISION,
    "enums" VARCHAR[],
    "min_int" BIGINT,
    "max_int" BIGINT,
    "min_float" DOUBLE PRECISION,
    "max_float" DOUBLE PRECISION,
    "float_precision" BIGINT,
    "order_index" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "game_args";
-- +goose StatementEnd
