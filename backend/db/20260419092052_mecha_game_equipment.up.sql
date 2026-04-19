-- Mecha equipment catalog and integration columns.
--
-- Equipment is strictly additive: it enhances existing chassis/weapon
-- capabilities and never gates a base behaviour. Equipment consumes the
-- same slot budget on the chassis as weapons, via the Mountable
-- abstraction in the domain layer.
--
-- Ammo now lives on the weapon itself (ammo_capacity > 0 means the weapon
-- has a built-in magazine). Ammo-bin equipment is purely additional
-- reloads layered on top. A weapon with ammo_capacity = 0 is an energy
-- weapon and never consumes from the pool.
BEGIN;

CREATE TABLE public.mecha_game_equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    mount_size VARCHAR(20) NOT NULL DEFAULT 'medium',
    effect_kind VARCHAR(30) NOT NULL,
    magnitude INTEGER NOT NULL DEFAULT 1,
    heat_cost INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_game_equipment_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_game_equipment_mount_size_check CHECK (mount_size IN ('small', 'medium', 'large')),
    CONSTRAINT mecha_game_equipment_effect_kind_check CHECK (effect_kind IN (
        'heat_sink', 'targeting_computer', 'armor_upgrade', 'jump_jets', 'ecm', 'ammo_bin'
    )),
    CONSTRAINT mecha_game_equipment_magnitude_check CHECK (magnitude BETWEEN 1 AND 200),
    CONSTRAINT mecha_game_equipment_heat_cost_check CHECK (heat_cost BETWEEN 0 AND 20)
);
CREATE INDEX idx_mecha_game_equipment_game_id ON public.mecha_game_equipment(game_id);
COMMENT ON TABLE public.mecha_game_equipment IS 'Designer-authored equipment catalog. Equipment enhances chassis/weapon capabilities and consumes chassis slots alongside weapons.';

ALTER TABLE public.mecha_game_weapon
    ADD COLUMN ammo_capacity INTEGER NOT NULL DEFAULT 0,
    ADD CONSTRAINT mecha_game_weapon_ammo_capacity_check CHECK (ammo_capacity BETWEEN 0 AND 200);

ALTER TABLE public.mecha_game_squad_mech
    ADD COLUMN equipment_config JSONB NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE public.mecha_game_mech_instance
    ADD COLUMN equipment_config JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN ammo_remaining INTEGER NOT NULL DEFAULT 0;

COMMIT;
