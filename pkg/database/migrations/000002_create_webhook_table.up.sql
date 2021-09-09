CREATE TABLE IF NOT EXISTS webhook(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    body TEXT NOT NULL
);