-- Dữ liệu demo
INSERT INTO wallets (admin_id, currency, balance) VALUES
    ('alice', 'VND', 5000000),
    ('bob',   'VND', 12000000)
ON CONFLICT (admin_id) DO NOTHING;

INSERT INTO transactions (id, admin_id, amount, note, created_by) VALUES
    ('tx-1', 'alice', 5000000,  'initial funding', 'system'),
    ('tx-1', 'bob',   12000000, 'initial funding', 'system')
ON CONFLICT (admin_id, id) DO NOTHING;
