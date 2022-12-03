SELECT
    user.firsname,
    user.lastname,
    purchase.sku,
    purchase.user_id
FROM
    user
    INNER JOIN purchase ON user.id = purchase.user_id
WHERE
    EXTRACT(
        MONTH
        FROM
            purchase.date
    ) = 2;