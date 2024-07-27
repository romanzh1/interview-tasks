-- +goose Up
-- +goose StatementBegin
create table if not exists items
(
    user_id bigint not null,
    sku     int not null,
    count   int not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items CASCADE;
-- +goose StatementEnd
