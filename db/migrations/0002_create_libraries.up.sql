CREATE TABLE IF NOT EXISTS libraries (
                                         id   BIGSERIAL PRIMARY KEY,
                                         name TEXT NOT NULL UNIQUE
);

ALTER TABLE books
    ADD COLUMN IF NOT EXISTS library_id BIGINT NULL;

ALTER TABLE books
    ADD CONSTRAINT fk_books_library
    FOREIGN KEY (library_id) REFERENCES libraries(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_books_library_id ON books (library_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_book_readings_triplet
    ON book_readings (book_id, user_id, read_date);