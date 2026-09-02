-- Dev/demo seed data so the app is testable end-to-end without manual entry.
-- Uses fixed UUIDs so the down migration can remove exactly these rows.

INSERT INTO categories (id, name, slug, description) VALUES
    ('11111111-1111-4111-8111-111111111101', 'Leafy Greens', 'leafy-greens', 'Spinach, methi, coriander and more'),
    ('11111111-1111-4111-8111-111111111102', 'Root Vegetables', 'root-vegetables', 'Potatoes, carrots, radish and more'),
    ('11111111-1111-4111-8111-111111111103', 'Exotic Vegetables', 'exotic-vegetables', 'Broccoli, zucchini, bell peppers and more'),
    ('11111111-1111-4111-8111-111111111104', 'Fresh Fruits', 'fresh-fruits', 'Seasonal fruits picked fresh');

INSERT INTO products (id, category_id, name, slug, description, unit, price, stock_quantity) VALUES
    ('22222222-2222-4222-8222-222222222201', '11111111-1111-4111-8111-111111111101', 'Spinach', 'spinach', 'Fresh green spinach', 'kg', 30.00, 100),
    ('22222222-2222-4222-8222-222222222202', '11111111-1111-4111-8111-111111111101', 'Coriander', 'coriander', 'Fresh coriander leaves', 'bundle', 10.00, 200),
    ('22222222-2222-4222-8222-222222222203', '11111111-1111-4111-8111-111111111102', 'Potato', 'potato', 'Farm fresh potatoes', 'kg', 25.00, 300),
    ('22222222-2222-4222-8222-222222222204', '11111111-1111-4111-8111-111111111102', 'Carrot', 'carrot', 'Crunchy orange carrots', 'kg', 40.00, 150),
    ('22222222-2222-4222-8222-222222222205', '11111111-1111-4111-8111-111111111103', 'Broccoli', 'broccoli', 'Fresh green broccoli', 'piece', 45.00, 60),
    ('22222222-2222-4222-8222-222222222206', '11111111-1111-4111-8111-111111111104', 'Banana', 'banana', 'Ripe bananas', 'dozen', 50.00, 80);

INSERT INTO delivery_areas (id, name, city, pincode) VALUES
    ('33333333-3333-4333-8333-333333333301', 'Sector 62', 'Noida', '201301'),
    ('33333333-3333-4333-8333-333333333302', 'Indirapuram', 'Ghaziabad', '201014'),
    ('33333333-3333-4333-8333-333333333303', 'Connaught Place', 'New Delhi', '110001');

INSERT INTO delivery_slots (id, label, start_time, end_time) VALUES
    ('44444444-4444-4444-8444-444444444401', 'Morning (6 AM - 9 AM)', '06:00', '09:00'),
    ('44444444-4444-4444-8444-444444444402', 'Afternoon (12 PM - 2 PM)', '12:00', '14:00'),
    ('44444444-4444-4444-8444-444444444403', 'Evening (5 PM - 8 PM)', '17:00', '20:00');
