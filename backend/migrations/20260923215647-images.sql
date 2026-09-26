
-- +migrate Up
CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT(uuidv7()),
    name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    seed BIGINT NOT NULL,
    style VARCHAR(32) NOT NULL,
    palette VARCHAR(32) NOT NULL,
    additional VARCHAR(32),
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    scale DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    uploaded BOOLEAN NOT NULL DEFAULT false
);

-- +migrate Down

DROP TABLE images;
