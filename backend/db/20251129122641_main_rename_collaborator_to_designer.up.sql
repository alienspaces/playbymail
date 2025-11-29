-- Rename Collaborator subscription type to Designer
-- This better reflects the purpose: designing games vs managing/playing them

-- Update existing records (if any)
UPDATE public.game_subscription
SET subscription_type = 'Designer'
WHERE subscription_type = 'Collaborator';

-- Drop and recreate the CHECK constraint with new value
ALTER TABLE public.game_subscription
DROP CONSTRAINT IF EXISTS game_subscription_subscription_type_check;

ALTER TABLE public.game_subscription
ADD CONSTRAINT game_subscription_subscription_type_check
CHECK (subscription_type IN ('Player', 'Manager', 'Designer'));

-- Update comment
COMMENT ON COLUMN public.game_subscription.subscription_type IS 'Role: Player, Manager, or Designer.';

