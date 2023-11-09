CREATE TABLE withdrawals
(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    create_at TIMESTAMPTZ NOT NULL,
    order_num BIGINT NOT NULL UNIQUE,
    sum INTEGER,
    user_id UUID NOT NULL,
    CONSTRAINT fk_customer
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE INDEX index_idx_withdrawals ON withdrawals (user_id);