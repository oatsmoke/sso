-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id        serial primary key not null,
    login     varchar unique     not null,
    pass_hash varchar            not null,
    email     varchar unique     not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
