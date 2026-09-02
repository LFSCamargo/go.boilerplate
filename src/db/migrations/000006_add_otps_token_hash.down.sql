DROP INDEX IF EXISTS idx_otps_token_hash;

ALTER TABLE otps
    DROP COLUMN IF EXISTS token_hash;
