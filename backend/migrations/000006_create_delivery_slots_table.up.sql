CREATE TABLE delivery_slots (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label      VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time   TIME NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (end_time > start_time)
);

CREATE TRIGGER delivery_slots_set_updated_at
    BEFORE UPDATE ON delivery_slots
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
