CREATE TABLE otp_codes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone       VARCHAR(15) NOT NULL,
    code_hash   TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    attempts    SMALLINT NOT NULL DEFAULT 0,
    consumed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Every lookup is "most recent row for this phone" (resend cooldown, hourly rate
-- limit, and verification all key off it), so that's the one index this needs.
CREATE INDEX otp_codes_phone_created_at_idx ON otp_codes (phone, created_at DESC);
