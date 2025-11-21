-- Create account_contact table and update game_subscription to reference it
-- Also remove name from account table

-- Create account_contact table
CREATE TABLE IF NOT EXISTS public.account_contact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES public.account(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    postal_address_line1 VARCHAR(255) NOT NULL,
    postal_address_line2 VARCHAR(255),
    state_province VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

COMMENT ON TABLE public.account_contact IS 'Contact information for accounts, including name and postal address. Collected from join game turn sheets.';
COMMENT ON COLUMN public.account_contact.id IS 'Unique identifier for the contact record.';
COMMENT ON COLUMN public.account_contact.account_id IS 'The account this contact belongs to.';
COMMENT ON COLUMN public.account_contact.name IS 'Contact name.';
COMMENT ON COLUMN public.account_contact.postal_address_line1 IS 'Primary postal address line.';
COMMENT ON COLUMN public.account_contact.postal_address_line2 IS 'Secondary postal address line (optional).';
COMMENT ON COLUMN public.account_contact.state_province IS 'State or province.';
COMMENT ON COLUMN public.account_contact.country IS 'Country.';
COMMENT ON COLUMN public.account_contact.postal_code IS 'Postal or ZIP code.';

-- Create indexes
CREATE INDEX idx_account_contact_account_id ON public.account_contact(account_id);

-- Add account_contact_id to game_subscription
ALTER TABLE public.game_subscription
    ADD COLUMN account_contact_id UUID REFERENCES public.account_contact(id) ON DELETE SET NULL;

-- Make account_contact_id required for new subscriptions
-- Note: Existing subscriptions will have NULL account_contact_id temporarily
-- We'll handle migration of existing data separately if needed

-- Update game_subscription comments
COMMENT ON COLUMN public.game_subscription.account_contact_id IS 'The contact information used for this subscription. References account_contact collected from join game turn sheets.';

-- Remove name from account table
ALTER TABLE public.account
    DROP COLUMN IF EXISTS name;

-- Update account comments
COMMENT ON TABLE public.account IS 'User accounts for authentication. Contact information is stored separately in account_contact.';

