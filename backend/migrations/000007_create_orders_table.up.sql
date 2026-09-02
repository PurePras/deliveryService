CREATE TABLE orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    delivery_area_id  UUID NOT NULL REFERENCES delivery_areas(id) ON DELETE RESTRICT,
    delivery_slot_id  UUID NOT NULL REFERENCES delivery_slots(id) ON DELETE RESTRICT,
    delivery_date     DATE NOT NULL,
    delivery_address  TEXT NOT NULL,
    status            VARCHAR(20) NOT NULL DEFAULT 'pending'
                          CHECK (status IN ('pending', 'confirmed', 'out_for_delivery', 'delivered', 'cancelled')),
    total_amount      NUMERIC(10, 2) NOT NULL CHECK (total_amount >= 0),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX orders_user_id_idx ON orders(user_id);
CREATE INDEX orders_area_date_slot_idx ON orders(delivery_area_id, delivery_date, delivery_slot_id);

CREATE TRIGGER orders_set_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
