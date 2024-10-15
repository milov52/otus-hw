CREATE table event (
                       id              UUID PRIMARY KEY,
                       title           text,
                       start_time      TIMESTAMP not null default now(),
                       duration        interval,
                       description     text,
                       user_id         text,
                       notify_before   interval,
                       sent            boolean default false,
                       created_at      TIMESTAMP not null default now(),
                       updated_at      DATE
);