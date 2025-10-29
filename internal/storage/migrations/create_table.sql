-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subs (
    name_service TEXT NOT NULL,
    cost_per_month INT DEFAULT 0 NOT NULL,
    user_id TEXT NOT NULL,
    date_start TEXT NOT NULL,
    date_end TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subs;
-- +goose StatementEnd