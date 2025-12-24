-- +goose Up
-- +goose StatementBegin
INSERT INTO "assets" ("id", "url", "size", )
VALUES
    ('f3d6c18b-2e47-4e52-9a3d-20d2a2e71a45', 'https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/3bc2e591-7f06-4273-b443-2d3cf6346c3c/d3jg07j-e2abb693-815f-4444-908c-fa94aab410a1.png?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7InBhdGgiOiJcL2ZcLzNiYzJlNTkxLTdmMDYtNDI3My1iNDQzLTJkM2NmNjM0NmMzY1wvZDNqZzA3ai1lMmFiYjY5My04MTVmLTQ0NDQtOTA4Yy1mYTk0YWFiNDEwYTEucG5nIn1dXSwiYXVkIjpbInVybjpzZXJ2aWNlOmZpbGUuZG93bmxvYWQiXX0.FWNFfuD5LpsbWOA9jEhLKckS_-3Gp-lP9vOqlXr8Sbk'),
    ('9e1b2b4d-6a3f-4f08-9a61-b74f8cf63d87', 'https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/b1f444f5-ca6c-4df3-8c44-ec6e0276f5ac/d7f0049-a816d6ab-034c-4c2e-8091-82b8486000ea.png/v1/fit/w_512,h_512/counter_strike_1_6_icon_by_dudekpro_d7f0049-375w-2x.png?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7ImhlaWdodCI6Ijw9NTEyIiwicGF0aCI6IlwvZlwvYjFmNDQ0ZjUtY2E2Yy00ZGYzLThjNDQtZWM2ZTAyNzZmNWFjXC9kN2YwMDQ5LWE4MTZkNmFiLTAzNGMtNGMyZS04MDkxLTgyYjg0ODYwMDBlYS5wbmciLCJ3aWR0aCI6Ijw9NTEyIn1dXSwiYXVkIjpbInVybjpzZXJ2aWNlOmltYWdlLm9wZXJhdGlvbnMiXX0.rcv0AOUaOMSe_QUK9XVIAv3ka786b7K3pRSL4dg8ts0')
ON CONFLICT DO NOTHING;

INSERT INTO "game_blobs" ("game_slug", "asset_id") 
VALUES 
    ('call-of-duty-2', 'f3d6c18b-2e47-4e52-9a3d-20d2a2e71a45'),
    ('counter-strike-1-6', '9e1b2b4d-6a3f-4f08-9a61-b74f8cf63d87')
ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "game_blobs";
-- +goose StatementEnd
