CREATE TABLE IF NOT EXISTS books (
                                     id     BIGSERIAL PRIMARY KEY,
                                     title  TEXT NOT NULL,
                                     author TEXT NOT NULL,
                                     year   INT  NOT NULL CHECK (year > 0),
    pages  INT  NOT NULL CHECK (pages > 0)
    );

CREATE TABLE IF NOT EXISTS users (
                                     id        BIGSERIAL PRIMARY KEY,
                                     full_name TEXT NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS book_readings (
    book_id   INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_date DATE NOT NULL,
    PRIMARY KEY (book_id, user_id)
    );

CREATE INDEX IF NOT EXISTS idx_books_title  ON books (title);
CREATE INDEX IF NOT EXISTS idx_books_author ON books (author);
CREATE INDEX IF NOT EXISTS idx_books_year   ON books (year);
CREATE INDEX IF NOT EXISTS idx_readings_book ON book_readings (book_id);

ALTER TABLE book_readings
    DROP CONSTRAINT book_readings_pkey;

ALTER TABLE book_readings
    ADD CONSTRAINT book_readings_pkey PRIMARY KEY (book_id, user_id, read_date);