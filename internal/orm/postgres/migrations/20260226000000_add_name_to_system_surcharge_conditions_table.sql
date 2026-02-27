-- +goose Up
-- +goose StatementBegin

ALTER TABLE system_surcharge_conditions
    ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE system_surcharge_conditions
    DROP COLUMN IF EXISTS name;

-- +goose StatementEnd

