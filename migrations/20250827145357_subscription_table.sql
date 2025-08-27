-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INT NOT NULL,
    user_id TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subscriptions;
-- +goose StatementEnd
