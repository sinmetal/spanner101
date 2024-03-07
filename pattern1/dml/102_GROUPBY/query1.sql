EXPLAIN ANALYZE
SELECT
    UserID,
    SUM(Orders.Amount) AS Amount
FROM
    Orders
GROUP BY
    UserID