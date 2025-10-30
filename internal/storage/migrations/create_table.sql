-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subs (
    id SERIAL NOT NULL PRIMARY KEY,
    name_service TEXT NOT NULL,
    cost_per_month INT DEFAULT 0 NOT NULL,
    user_uuid TEXT NOT NULL,
    date_start DATE NOT NULL,
    date_end DATE NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subs;
-- +goose StatementEnd