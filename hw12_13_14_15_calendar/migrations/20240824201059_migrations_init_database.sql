-- +goose Up
CREATE table event (
                       id              UUID PRIMARY KEY,
                       title           text,
                       start_time      TIMESTAMP not null default now(),
                       duration        timestamptz,
                       description     text,
                       user_id         text,
                       notification    timestamptz,
                       created_at      TIMESTAMP not null default now(),
                       updated_at      DATE
);

-- +goose Down
drop table event;