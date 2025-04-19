begin;
CREATE TABLE IF NOT EXISTS request (
    id UUID PRIMARY KEY, 
    result JSONB,
    status VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP  SET DEFAULT now()
);

CREATE TABLE outbox (
    id UUID PRIMARY KEY, 
    data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP  NOT NULL DEFAULT now()
);
commit;