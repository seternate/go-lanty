-- +goose Up
-- +goose StatementBegin
INSERT INTO "games" ("slug", "name")
VALUES 
    ('call-of-duty-2', 'Call of Duty 2'),
    ('rune', 'Rune'),
    ('counter-strike-1-6', 'Counter-Strike 1.6')
ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "games";
-- +goose StatementEnd
