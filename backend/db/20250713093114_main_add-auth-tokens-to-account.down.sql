ALTER TABLE account
  DROP COLUMN IF EXISTS verification_token,
  DROP COLUMN IF EXISTS verification_token_expires_at,
  DROP COLUMN IF EXISTS session_token,
  DROP COLUMN IF EXISTS session_token_expires_at;
