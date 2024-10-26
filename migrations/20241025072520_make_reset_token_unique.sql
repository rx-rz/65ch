-- +goose Up
-- +goose StatementBegin
ALTER TABLE reset_tokens
ADD CONSTRAINT unique_reset_token UNIQUE (reset_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reset_tokens
DROP CONSTRAINT unique_reset_token
-- +goose StatementEnd
