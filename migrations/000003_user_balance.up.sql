CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists user_balance (
    user_id uuid references users(id) on delete cascade ,
    current money default 0,
    withdrawn money default 0
);

create table if not exists withdrawals (
    user_id uuid references users(id) on delete cascade,
    order_number text unique,
    sum money not null,
    processed_at timestamp default now()
);
