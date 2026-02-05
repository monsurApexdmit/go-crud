-- Seed the 14 permission modules.
-- Uses INSERT IGNORE so re-running the seeder is safe (idempotent).
INSERT IGNORE INTO permissions (name) VALUES
('Dashboard'),
('Products'),
('Categories'),
('Attributes'),
('Coupons'),
('Customers'),
('Orders'),
('POS'),
('Sells'),
('Staff'),
('Settings'),
('International'),
('Store'),
('Pages');
