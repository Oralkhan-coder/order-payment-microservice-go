CREATE TYPE order_status AS ENUM ('pending', 'paid', 'failed', 'cancelled');


CREATE TABLE IF NOT EXISTS orders
(
    id          UUID PRIMARY KEY,
    customer_id UUID         NOT NULL,
    item_name   VARCHAR(255) NOT NULL,
    amount      BIGINT       NOT NULL,
    status      order_status NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMP    NOT NULL
);
