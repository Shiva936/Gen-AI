-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS document_embeddings(
    id BIGSERIAL PRIMARY KEY,
    document_id uuid,
    content TEXT,
    embedding VECTOR(768)
);
DROP CONSTRAINT IF EXISTS dockument_id_fkey;
AFTER TABLE document_embeddings ADD CONSTRAINT dockument_id_fkey FOREIGN KEY (document_id) REFERENCES documents (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP CONSTRAINT IF EXISTS dockument_id_fkey;
DROP TABLE IF EXISTS document_embeddings;
-- +goose StatementEnd
