-- 003_update_files_table_to_match_model.up.sql

ALTER TABLE files DROP COLUMN permissions;
ALTER TABLE files ALTER COLUMN lock_info TYPE TEXT;
ALTER TABLE files ALTER COLUMN lock_info DROP NOT NULL;
