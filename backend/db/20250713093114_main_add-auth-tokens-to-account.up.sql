ALTER TABLE account
  ADD COLUMN verification_token TEXT,
  ADD COLUMN verification_token_expires_at TIMESTAMPTZ,
  ADD COLUMN session_token TEXT,
  ADD COLUMN session_token_expires_at TIMESTAMPTZ;
