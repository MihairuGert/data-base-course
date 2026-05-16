CREATE OR REPLACE TABLE app_users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash CHAR(64) NOT NULL,
    role VARCHAR(30) NOT NULL CHECK (role IN (
        'management',
        'workshop_heads',
        'foremen',
        'dispatchers',
        'accounting',
        'hr',
        'drivers_role',
        'repairmen_role'
    )),
    employee_id BIGINT REFERENCES employees(id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO app_users (username, password_hash, role, employee_id) VALUES
('director', '710815eba69d3c488d3affd899583414d72f6aa811c9510312c1bb4ba7add839', 'management', 1),
('workshop', '89e8e6faca85a3ae6784c900c74ae8a05eed6a685ef6ac3ff9e58ffa738a622d', 'workshop_heads', 1),
('foreman', '743e1ad97d5e2c1f05d1bf9df3659cc3e1ae057422c36d350d3d6e17dc08d0a6', 'foremen', 2),
('dispatcher', '6f33e5ff723e36a0c178758087ff021ad7e91895a977a7c8ccfefeb0a3609d54', 'dispatchers', NULL),
('accountant', '45aa6c86b5ed1ab56d5157b117a3988654d50561f45e433c97356df9b2c22a9d', 'accounting', NULL),
('hr', '365906c37b8a3a6b82c6b9d1c4529e4cbf0fc666cc077feab4bf5ad03dd3580e', 'hr', NULL),
('driver', '009e26c03dabb8cc36d1e1949e0593594477aadc0cdc847a1209503f1d9bf523', 'drivers_role', 3),
('repairman', '0cbfae9dd3124e43a196b83763cb43bf8917ea2bf511a1e680b0b1d2410315de', 'repairmen_role', 4);
