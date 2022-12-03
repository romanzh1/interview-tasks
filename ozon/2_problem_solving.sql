CREATE TABLE userr (
    id SERIAL PRIMARY KEY,
    first VARCHAR(255),
    second VARCHAR(255)
);

INSERT INTO
    userr (first, second)
VALUES
    ('f1', 's1'),
    ('s2', 's2');

CREATE TABLE purchase (
    sku VARCHAR(255),
    price INTEGER,
    user_id INTEGER REFERENCES userr(id),
    date TIMESTAMP
);

DROP TABLE purchase;

INSERT INTO
    purchase (sku, price, user_id, date)
VALUES
    ('11111', 5500, 1, '2021-02-15'),
    ('22222', 4000, 1, '2021-02-14'),
    ('33333', 8000, 2, '2021-03-01'),
    ('44444', 400, 2, '2021-03-02');

CREATE TABLE ban_list (
    user_id INTEGER REFERENCES userr(id),
    date TIMESTAMP
);

INSERT INTO
    ban_list (user_id, date)
VALUES
    (1, '2021-03-08');

SELECT
    *
FROM
    userr;

SELECT
    userr.first,
    userr.second,
    SUM(purchase.price) as price
FROM
    userr
    INNER JOIN purchase ON userr.id = purchase.user_id
    LEFT JOIN ban_list ON userr.id = ban_list.user_id
WHERE
    ban_list.user_id IS NULL
GROUP BY
    price,
    userr.first,
    userr.second
HAVING
    price > 5000;