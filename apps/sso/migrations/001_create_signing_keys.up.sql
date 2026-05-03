CREATE TABLE IF NOT EXISTS signing_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    algorithm VARCHAR(16) NOT NULL DEFAULT 'RS256',
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    rotated_at TIMESTAMP
)
CREATE INDEX IF NOT EXISTS idx_signing_keys_active ON signing_keys(is_active) WHERE is_active = true;
