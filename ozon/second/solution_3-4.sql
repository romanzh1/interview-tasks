Автор
ID
    Name

Автор-книга
    ID_book
    ID_author

Книга
    ID
    ID_reader
    Name

Читатель
    ID
    Name

-- task 4

SELECT name, COUNT(ID_author) as auth
FROM Книга AS k
         INNER JOIN Автор-книга AS ab ON ab.ID_author = k.ID
GROUP BY name
HAVING COUNT(ID_author) > 3

