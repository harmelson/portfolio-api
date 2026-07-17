CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    google_id TEXT NOT NULL UNIQUE,
    picture_id TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);

CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    max_apps INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    purchase_token TEXT NOT NULL UNIQUE,
    order_id TEXT,
    status TEXT NOT NULL,
    auto_renew BOOLEAN NOT NULL DEFAULT TRUE,
    started_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_status_expires_at ON subscriptions(status, expires_at);

-- Seed plans
INSERT INTO plans (name, max_apps, price) VALUES
    ('free', 1, 0.00),
    ('basic', 1, 5.99),
    ('plus', -1, 8.99)
ON CONFLICT (name) DO NOTHING;

-- Seed test user
INSERT INTO users (id, email, name, google_id, picture_id, is_active)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'dev@example.com',
    'Test User',
    'dev-test-user',
    'https://example.com/picture.png',
    true
)
ON CONFLICT (google_id) DO NOTHING;

-- Seed free subscription for test user
INSERT INTO subscriptions (user_id, plan_id, purchase_token, order_id, status, auto_renew, started_at, expires_at)
SELECT
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    p.id,
    'free_a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    '',
    'free',
    false,
    NOW(),
    NOW() + INTERVAL '100 years'
FROM plans p
WHERE p.name = 'free'
ON CONFLICT (purchase_token) DO NOTHING;
