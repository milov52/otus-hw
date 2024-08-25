-- +goose Up
CREATE table event (
                       id              serial primary key,
                       title           text,
                       start_time      timestamptz not null default now() unique,
                       duration        timestamptz,
                       description     text,
                       user_id         text,
                       notification    timestamptz,
                       created_at      timestamptz not null default now(),
                       updated_at      timestamptz
);

-- +goose Down
drop table event;