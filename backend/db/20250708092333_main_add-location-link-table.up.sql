CREATE TABLE public.location_link (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_location_id UUID NOT NULL REFERENCES public.location(id),
    to_location_id UUID NOT NULL REFERENCES public.location(id),
    description VARCHAR(255) NOT NULL,
    name VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_location_link_from_location_id ON public.location_link(from_location_id);
CREATE INDEX idx_location_link_to_location_id ON public.location_link(to_location_id);

-- Prevent duplicate links
ALTER TABLE public.location_link ADD CONSTRAINT location_link_unique UNIQUE (from_location_id, to_location_id);
