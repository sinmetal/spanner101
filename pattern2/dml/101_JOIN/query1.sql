EXPLAIN ANALYZE
SELECT *
FROM Singers s INNER JOIN Albums a ON s.SingerId = a.SingerId
WHERE s.SingerId = 1;