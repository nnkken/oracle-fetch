-- +goose Up
-- +goose StatementBegin
ALTER TABLE prices ADD COLUMN source TEXT;
UPDATE prices SET source = 'chainlink-eth';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE prices DROP COLUMN source;
-- +goose StatementEnd
