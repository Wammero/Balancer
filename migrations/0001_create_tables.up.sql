-- Таблица токенов для айпишников
CREATE TABLE token_buckets (
    key VARCHAR(255) PRIMARY KEY,
    tokens INT NOT NULL,
    last_refill TIMESTAMP NOT NULL,
    tokens_per_second INT NOT NULL,
    capacity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_token_buckets_key ON token_buckets(key);
