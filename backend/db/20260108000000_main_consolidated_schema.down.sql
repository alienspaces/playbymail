-- Consolidated Schema Teardown
-- Generated: 2026-01-08

BEGIN;

DROP VIEW IF EXISTS public.game_subscription_view;


DROP TABLE IF EXISTS public.game_turn_sheet;
DROP TABLE IF EXISTS public.game_instance_parameter;
DROP TABLE IF EXISTS public.game_subscription_instance;
DROP TABLE IF EXISTS public.game_subscription;
DROP TABLE IF EXISTS public.game_instance;
DROP TABLE IF EXISTS public.game_image;
DROP TABLE IF EXISTS public.game;

DROP TABLE IF EXISTS public.account_subscription;
DROP TABLE IF EXISTS public.account_user_contact;
DROP TABLE IF EXISTS public.account_user;
DROP TABLE IF EXISTS public.account;

COMMIT;
