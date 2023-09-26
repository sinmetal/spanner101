# GROUP BY pattern2

```
EXPLAIN ANALYZE
SELECT
  UserID,
  SUM(Orders.Amount) AS Amount
FROM
  Orders
GROUP BY
  UserID
```

```
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                         | Rows_Returned | Executions | Total_Latency |
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)    | 5             | 1          | 3.15 msecs    |
|  1 | +- Serialize Result                                                          | 5             | 1          | 3.13 msecs    |
|  2 |    +- Stream Aggregate                                                       | 5             | 1          | 3.13 msecs    |
|  3 |       +- Local Distributed Union                                             | 4353          | 1          | 2.86 msecs    |
|  4 |          +- Table Scan (Full scan: true, Table: Orders, scan_method: Scalar) | 4353          | 1          | 2.6 msecs     |
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
5 rows in set (7.53 msecs)
timestamp:            2023-09-15T14:35:53.678128+09:00
cpu time:             6.77 msecs
rows scanned:         4353 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230906_07_18_51UTC
```

PKの先頭にUserIDがあるので、UserごとにOrderはまとまっている。
そのため、Stream Aggregateで集計ができる。
Orders TableはUsers Tableの子どもなので、Localで完結している。