-- Дано:
CREATE TABLE cities
(
    id   serial PRIMARY KEY,
    name text NOT NULL
);

INSERT INTO cities (name)
VALUES ('Москва'),
       ('Санкт-Петербург'),
       ('Краснодар');

CREATE TABLE users
(
    id      serial PRIMARY KEY,
    name    text NOT NULL,
    city_id int  NOT NULL REFERENCES cities (id)
);

INSERT INTO users (name, city_id)
VALUES ('Иван', 1),
       ('Анна', 1),
       ('Олег', 2);

-- Решение сложных кейсов:
-- Города с пользователями и без
SELECT cities.name, users.name
FROM cities
         LEFT JOIN users ON users.city_id = cities.id;

-- Города с пользователями и без с сортировкой по возрастанию и номеру строки по убыванию
SELECT cities.name, NULLIF(COUNT(users.id), 0) AS users_count, @row_number() OVER (ORDER BY COUNT(users.id) DESC)
FROM cities
         LEFT JOIN users ON users.city_id = cities.id
GROUP BY cities.name
ORDER BY COUNT(users.id);

