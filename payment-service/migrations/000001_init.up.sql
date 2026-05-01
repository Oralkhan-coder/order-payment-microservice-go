CREATE TYPE payment_status AS ENUM ('Authorized', 'Declined');


CREATE TABLE IF NOT EXISTS payments
(
    id             UUID PRIMARY KEY,
    order_id       UUID           NOT NULL,
    transaction_id UUID           NOT NULL,
    amount         BIGINT         NOT NULL,
    status         payment_status NOT NULL DEFAULT 'Declined'
);
