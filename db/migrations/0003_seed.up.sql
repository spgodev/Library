ALTER TABLE users
    ADD COLUMN library_id BIGINT;

UPDATE users
SET library_id = (SELECT id FROM libraries WHERE name = 'Default Library')
WHERE library_id IS NULL;

ALTER TABLE users
    ALTER COLUMN library_id SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT fk_users_library
        FOREIGN KEY (library_id) REFERENCES libraries(id) ON DELETE RESTRICT;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_full_name_key;

CREATE UNIQUE INDEX uq_users_library_full_name
    ON users (library_id, full_name);

CREATE INDEX idx_users_library_id
    ON users (library_id);