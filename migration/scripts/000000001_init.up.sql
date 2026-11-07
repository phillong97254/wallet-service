CREATE TABLE IF NOT EXISTS wallets (
    admin_id   TEXT PRIMARY KEY,
    currency   TEXT        NOT NULL DEFAULT 'VND',
    balance    BIGINT      NOT NULL DEFAULT 0 CHECK (balance >= 0), -- đơn vị nhỏ nhất (đồng)
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transactions (
    seq        BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, -- thứ tự ghi
    id         TEXT        NOT NULL,                             -- tx-1, tx-2... theo từng ví
    admin_id   TEXT        NOT NULL REFERENCES wallets (admin_id),
    amount     BIGINT      NOT NULL,                             -- âm = trừ tiền
    note       TEXT        NOT NULL DEFAULT '',
    created_by TEXT        NOT NULL,                             -- audit: ai đã sửa
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (admin_id, id)
);
