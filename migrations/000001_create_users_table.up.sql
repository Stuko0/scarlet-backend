SET search_path TO scarlet;
create table users(
    user_id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT UNIQUE,
    password TEXT,
    phone TEXT,
    origin TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    image TEXT,
    role TEXT NOT NULL DEFAULT 'civil',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);