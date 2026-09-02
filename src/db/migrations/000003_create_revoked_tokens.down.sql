DROP INDEX IF EXISTS idx_revoked_tokens_expires_at;
DROP INDEX IF EXISTS idx_revoked_tokens_user_id;
DROP INDEX IF EXISTS idx_revoked_tokens_jti;
DROP TABLE IF EXISTS revoked_tokens;
