INSERT INTO books (id, title, author, year, pages) VALUES
                                                       (1, 'Grokking Algorithms',     'Aditya Bhargava',          2016, 278),
                                                       (2, 'Design Patterns',         'Head First',               2022, 633),
                                                       (3, '1984',                    'George Orwell',            1949, 210),
                                                       (4, 'The Little Prince',       'Antoine de Saint-Exupery', 1943, 60),
                                                       (5, 'Crime and Punishment',    'Fyodor Dostoevsky',        1886, 321)
    ON CONFLICT (id) DO NOTHING;

INSERT INTO users (full_name) VALUES
                                  ('Anton'),
                                  ('Sergey'),
                                  ('Vasya Pupkin')
    ON CONFLICT (full_name) DO NOTHING;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT 1, u.id, DATE '2020-03-10' FROM users u WHERE u.full_name = 'Anton'
    ON CONFLICT (book_id, user_id) DO UPDATE SET read_date = EXCLUDED.read_date;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT 1, u.id, DATE '2025-12-10' FROM users u WHERE u.full_name = 'Sergey'
    ON CONFLICT (book_id, user_id) DO UPDATE SET read_date = EXCLUDED.read_date;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT 2, u.id, DATE '2025-12-25' FROM users u WHERE u.full_name = 'Sergey'
    ON CONFLICT (book_id, user_id) DO UPDATE SET read_date = EXCLUDED.read_date;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT 5, u.id, DATE '2007-10-30' FROM users u WHERE u.full_name = 'Vasya Pupkin'
    ON CONFLICT (book_id, user_id) DO UPDATE SET read_date = EXCLUDED.read_date;