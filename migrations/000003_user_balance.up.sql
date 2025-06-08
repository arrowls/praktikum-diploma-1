create table if not exists user_balance (
    user_id uuid references users(id),
    current float default 0,
    withdrawn float default 0
);

create table if not exists withdrawals (
    user_id uuid references users(id),
    order_number text unique,
    sum float not null,
    processed_at timestamp default now()
);
