CREATE TABLE IF NOT EXISTS public.account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL DEFAULT '',
    verification_token TEXT,
    verification_token_expires_at TIMESTAMPTZ,
    session_token TEXT,
    session_token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE
);
