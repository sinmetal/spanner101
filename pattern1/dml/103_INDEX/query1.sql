EXPLAIN ANALYZE
SELECT
    UserID,
    OrderID,
    Amount,
    CommitedAt
FROM Orders
WHERE UserID = "ruby"
ORDER BY CommitedAt DESC
LIMIT 5