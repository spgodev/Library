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
                                             id        BIGSERIAL PRIMARY KEY,
                                             book_id   BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
                                             user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                             read_date DATE   NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_books_title   ON books (title);
CREATE INDEX IF NOT EXISTS idx_books_author  ON books (author);
CREATE INDEX IF NOT EXISTS idx_books_year    ON books (year);

CREATE INDEX IF NOT EXISTS idx_readings_book ON book_readings (book_id);
CREATE INDEX IF NOT EXISTS idx_readings_user ON book_readings (user_id);
CREATE INDEX IF NOT EXISTS idx_readings_date ON book_readings (read_date);