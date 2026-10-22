-- +goose Up
CREATE TABLE examples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (char_length(name) BETWEEN 1 AND 120)
);

-- +goose Down
DROP TABLE examples;
