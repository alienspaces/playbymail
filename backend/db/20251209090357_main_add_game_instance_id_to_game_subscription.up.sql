-- Add game_instance_id to game_subscription table
-- This field links Player subscriptions to the specific game instance they are playing in
-- Only valid for Player type subscriptions

ALTER TABLE public.game_subscription
ADD COLUMN game_instance_id UUID;

-- Add foreign key constraint
ALTER TABLE public.game_subscription
ADD CONSTRAINT game_subscription_game_instance_id_fkey 
FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id);

-- Add index for efficient lookups
CREATE INDEX idx_game_subscription_game_instance_id ON public.game_subscription(game_instance_id)
WHERE game_instance_id IS NOT NULL;

-- Add comment
COMMENT ON COLUMN public.game_subscription.game_instance_id IS 'The game instance this Player subscription is associated with. Only valid for Player type subscriptions.';

