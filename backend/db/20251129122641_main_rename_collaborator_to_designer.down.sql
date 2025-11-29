-- Revert Designer subscription type back to Collaborator

-- Update existing records
UPDATE public.game_subscription
SET subscription_type = 'Collaborator'
WHERE subscription_type = 'Designer';

-- Drop and recreate the CHECK constraint with original value
ALTER TABLE public.game_subscription
DROP CONSTRAINT IF EXISTS game_subscription_subscription_type_check;

ALTER TABLE public.game_subscription
ADD CONSTRAINT game_subscription_subscription_type_check
CHECK (subscription_type IN ('Player', 'Manager', 'Collaborator'));

-- Revert comment
COMMENT ON COLUMN public.game_subscription.subscription_type IS 'Role: Player, Manager, or Collaborator.';

