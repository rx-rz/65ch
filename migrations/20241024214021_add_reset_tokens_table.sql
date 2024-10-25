-- +goose Up
-- +goose StatementBegin
CREATE TABLE reset_tokens(
    id serial primary key ,
    user_id uuid references users(id),
    reset_token text not null,
    expiration timestamptz not null,
    is_used boolean default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reset_tokens;
-- +goose StatementEnd
