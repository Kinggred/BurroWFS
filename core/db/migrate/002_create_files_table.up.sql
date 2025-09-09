-- 002_create_files_table.up.sql
CREATE TABLE files (
    id UUID PRIMARY KEY,
    file_id UUID ,                     -- Logical file identifier, null if folder
    version_parent_id UUID REFERENCES files(id) DEFERRABLE INITIALLY DEFERRED , -- previous version
    owner_id UUID NOT NULL REFERENCES users(id),
    path_parent_id UUID REFERENCES files(id) DEFERRABLE INITIALLY DEFERRED ,  -- Folder hierarchy
    name TEXT NOT NULL,
    path TEXT NOT NULL,                       -- Full path for easy lookup
    s3_key TEXT NOT NULL,
    size BIGINT,
    content_type TEXT,
    etag TEXT,
    version INT NOT NULL,                       -- Latest version
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    permissions JSONB,
    lock_info JSONB
);

CREATE INDEX idx_files_path_parent ON files(path_parent_id);

CREATE INDEX idx_files_fileid_version ON files(file_id, version DESC);

CREATE INDEX idx_files_folder_order ON files(path_parent_id, name);

CREATE OR REPLACE FUNCTION update_modified_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to files table
CREATE TRIGGER trg_update_files_modified
    BEFORE UPDATE ON files
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at();



CREATE TABLE file_versions (
   id UUID PRIMARY KEY,
   file_id UUID NOT NULL,                     -- Links to main file
   version_parent_id UUID REFERENCES file_versions(id), -- previous version
   owner_id UUID NOT NULL REFERENCES users(id),
   path_parent_id UUID,                       -- Folder hierarchy at the time
   name TEXT NOT NULL,
   s3_key TEXT NOT NULL,
   size BIGINT,
   content_type TEXT,
   is_folder BOOLEAN DEFAULT FALSE,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP DEFAULT NOW(),
   permissions JSONB,
   lock_info JSONB
);


CREATE INDEX idx_file_versions_path_parent ON file_versions(path_parent_id);