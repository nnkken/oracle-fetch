-- +goose Up
-- +goose StatementBegin
CREATE TABLE prices (
  id SERIAL PRIMARY KEY,
  token TEXT,
  unit TEXT,
  price NUMERIC,
  price_timestamp TIMESTAMP WITH TIME ZONE,
  fetch_timestamp TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE prices;
-- +goose StatementEnd
