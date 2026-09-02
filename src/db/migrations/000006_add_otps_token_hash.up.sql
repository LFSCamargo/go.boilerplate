ALTER TABLE otps
    ADD COLUMN IF NOT EXISTS token_hash VARCHAR(255);

CREATE UNIQUE INDEX IF NOT EXISTS idx_otps_token_hash
    ON otps (token_hash)
    WHERE token_hash IS NOT NULL;
