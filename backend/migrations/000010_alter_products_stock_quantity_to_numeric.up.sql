ALTER TABLE products
    ALTER COLUMN stock_quantity TYPE NUMERIC(10, 2) USING stock_quantity::NUMERIC(10, 2);
