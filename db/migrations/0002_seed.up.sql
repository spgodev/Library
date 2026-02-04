INSERT INTO books (title, author, year, pages) VALUES
                                                   ('Grokking Algorithms',  'Aditya Bhargava',          2016, 278),
                                                   ('Design Patterns',      'Head First',               2022, 633),
                                                   ('1984',                 'George Orwell',            1949, 210),
                                                   ('The Little Prince',    'Antoine de Saint-Exupery', 1943, 60),
                                                   ('Crime and Punishment', 'Fyodor Dostoevsky',        1886, 321)
ON CONFLICT DO NOTHING;

INSERT INTO users (full_name) VALUES
                                  ('Anton'),
                                  ('Sergey'),
                                  ('Vasya Pupkin')
ON CONFLICT (full_name) DO NOTHING;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT b.id, u.id, DATE '2020-03-10'
FROM books b
         JOIN users u ON u.full_name = 'Anton'
WHERE b.title = 'Grokking Algorithms'
ON CONFLICT (book_id, user_id, read_date) DO NOTHING;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT b.id, u.id, DATE '2025-12-10'
FROM books b
         JOIN users u ON u.full_name = 'Sergey'
WHERE b.title = 'Grokking Algorithms'
ON CONFLICT (book_id, user_id, read_date) DO NOTHING;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT b.id, u.id, DATE '2025-12-25'
FROM books b
         JOIN users u ON u.full_name = 'Sergey'
WHERE b.title = 'Design Patterns'
ON CONFLICT (book_id, user_id, read_date) DO NOTHING;

INSERT INTO book_readings (book_id, user_id, read_date)
SELECT b.id, u.id, DATE '2007-10-30'
FROM books b
         JOIN users u ON u.full_name = 'Vasya Pupkin'
WHERE b.title = 'Crime and Punishment'
ON CONFLICT (book_id, user_id, read_date) DO NOTHING;