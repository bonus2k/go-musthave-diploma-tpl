CREATE TABLE orders
(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    create_at TIMESTAMPTZ NOT NULL,
    number BIGINT NOT NULL UNIQUE,
    accrual DECIMAL,
    status VARCHAR(15) NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_customer
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE INDEX index_idx_orders ON orders (user_id, number, status);