CREATE TABLE user_cars (
    user_id INTEGER,
    car_id INTEGER NOT NULL
);

CREATE TABLE users (id INTEGER NOT NULL, name VARCHAR(255));

INSERT INTO
    user_cars (user_id, car_id)
VALUES
    (1, 5),
    (2, 1),
    (3, 1),
    (4, 7);

INSERT INTO
    purchase (sku, price, user_id, date)
VALUES
    ('11111', 5500, 1, '2021-02-15'),
    ('22222', 4000, 1, '2021-02-14'),
    ('33333', 8000, 2, '2021-03-01'),
    ('44444', 400, 2, '2021-03-02');