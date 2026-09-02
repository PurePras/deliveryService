CREATE TABLE delivery_areas (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(150) NOT NULL,
    city       VARCHAR(100) NOT NULL,
    pincode    VARCHAR(10) NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (pincode, name)
);

CREATE TRIGGER delivery_areas_set_updated_at
    BEFORE UPDATE ON delivery_areas
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
