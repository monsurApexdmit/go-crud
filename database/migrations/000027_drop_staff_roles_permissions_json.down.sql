ALTER TABLE staff_roles ADD COLUMN permissions JSON NOT NULL DEFAULT ('[]');
