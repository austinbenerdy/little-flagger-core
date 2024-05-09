-- +goose Up
-- +goose StatementBegin
CREATE TABLE feature_flag
(
    id          varchar(255) not null,
    name        varchar(255) not null,
    slug        varchar(255) not null,
    description varchar(255) not null,
    status      varchar(255) not null,
    percentage  int
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feature_flag
-- +goose StatementEnd
