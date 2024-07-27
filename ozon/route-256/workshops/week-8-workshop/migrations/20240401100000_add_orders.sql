-- +goose Up
-- +goose StatementBegin

CREATE TABLE orders
(
    id          BIGSERIAL NOT NULL PRIMARY KEY,
    user_id     BIGINT    NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE orders;

-- +goose StatementEnd
