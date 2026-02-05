ALTER TABLE vendors
    DROP FOREIGN KEY fk_vendors_user,
    DROP COLUMN user_id;
