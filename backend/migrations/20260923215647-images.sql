
-- +migrate Up
CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT(uuidv7()),
    name VARCHAR(255) NOT NULL,
    key VARCHAR(1024) NOT NULL UNIQUE,
    size BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +migrate Down

DROP TABLE images;
