-- +goose Up
-- +goose StatementBegin

CREATE SEQUENCE order_id_manual_seq INCREMENT 1000 START 1000; -- 1000 > number of buckets

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP SEQUENCE order_id_manual_seq;

-- +goose StatementEnd
