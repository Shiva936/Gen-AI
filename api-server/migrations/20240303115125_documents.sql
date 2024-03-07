-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS documents (
    id uuid PRIMARY KEY,
    name varchar
    url text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
