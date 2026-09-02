CREATE TABLE IF NOT EXISTS user_registration_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    require_email_verification BOOLEAN NOT NULL DEFAULT TRUE,
    min_password_length INT NOT NULL DEFAULT 8,
    allow_registration BOOLEAN NOT NULL DEFAULT TRUE,
    otp_expiry_minutes INT NOT NULL DEFAULT 10,
    max_otp_attempts INT NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_min_password_length CHECK (min_password_length >= 6),
    CONSTRAINT chk_otp_expiry_minutes CHECK (otp_expiry_minutes >= 1),
    CONSTRAINT chk_max_otp_attempts CHECK (max_otp_attempts >= 1)
);

INSERT INTO user_registration_configs (id)
SELECT gen_random_uuid()
WHERE NOT EXISTS (SELECT 1 FROM user_registration_configs);
