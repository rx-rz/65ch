-- +goose Up
-- +goose StatementBegin
ALTER TABLE reset_tokens
    DROP COLUMN is_used;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reset_tokens
    ADD COLUMN is_used BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd
