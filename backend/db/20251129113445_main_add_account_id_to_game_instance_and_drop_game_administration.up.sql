-- Add game_subscription_id to game_instance for ownership/RLS
-- This links instances to the subscription that created them, enabling future
-- organisation-based access control when organisations are introduced.

-- First, add the column as nullable
ALTER TABLE public.game_instance
ADD COLUMN game_subscription_id UUID REFERENCES public.game_subscription(id);

-- Update existing game instances to use the Manager subscription for their game
-- If no Manager subscription exists, use the first subscription for that game
UPDATE public.game_instance gi
SET game_subscription_id = (
    SELECT gs.id
    FROM public.game_subscription gs
    WHERE gs.game_id = gi.game_id
    AND gs.deleted_at IS NULL
    ORDER BY 
        CASE WHEN gs.subscription_type = 'Manager' THEN 0 ELSE 1 END,
        gs.created_at ASC
    LIMIT 1
);

-- For any remaining instances without a subscription (orphaned), delete them
DELETE FROM public.game_instance
WHERE game_subscription_id IS NULL;

-- Now make the column NOT NULL
ALTER TABLE public.game_instance
ALTER COLUMN game_subscription_id SET NOT NULL;

-- Add index for RLS queries
CREATE INDEX idx_game_instance_game_subscription_id ON public.game_instance(game_subscription_id);

-- Add comment
COMMENT ON COLUMN public.game_instance.game_subscription_id IS 'The Manager subscription that created this game instance.';

-- Drop game_administration table (replaced by Manager subscription type)
DROP TABLE IF EXISTS public.game_administration;

