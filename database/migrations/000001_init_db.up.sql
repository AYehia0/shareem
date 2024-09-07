CREATE table share (
    id uuid primary key,
    url varchar(256) not null,
    title varchar(256) not null,
    note varchar(256),
    ip inet not null,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

CREATE index ON share (created_at);
