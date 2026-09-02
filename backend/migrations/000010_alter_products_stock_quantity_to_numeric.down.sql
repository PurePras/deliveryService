ALTER TABLE products
    ALTER COLUMN stock_quantity TYPE INTEGER USING stock_quantity::INTEGER;
