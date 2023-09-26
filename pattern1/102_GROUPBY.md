# GROUP BY pattern1

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
+----+--------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                 | Rows_Returned | Executions | Total_Latency |
+----+--------------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Serialize Result                                                                     | 5             | 1          | 6.2 msecs     |
|  1 | +- Global Stream Aggregate                                                           | 5             | 1          | 6.2 msecs     |
|  2 |    +- Sort                                                                           | 5             | 1          | 6.19 msecs    |
|  3 |       +- Distributed Union (distribution_table: Orders, split_ranges_aligned: false) | 5             | 1          | 6.18 msecs    |
|  4 |          +- Local Hash Aggregate                                                     | 5             | 1          | 6.15 msecs    |
|  5 |             +- Local Distributed Union                                               | 4419          | 1          | 5.63 msecs    |
|  6 |                +- Table Scan (Full scan: true, Table: Orders, scan_method: Scalar)   | 4419          | 1          | 5.37 msecs    |
+----+--------------------------------------------------------------------------------------+---------------+------------+---------------+
5 rows in set (10.55 msecs)
timestamp:            2023-09-15T14:38:43.309677+09:00
cpu time:             8.39 msecs
rows scanned:         4419 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230909_06_42_57UTC
```

GROUP BYのために `Local Hash Aggregate` が登場。
この例だとUserIDの種類が5つしかないので、すぐ終わっているが、Userが多いとHash Table作って集計するのが大変になる。
Localでの集計が終わった後、 `Sort` して、 `Global Stream Aggregate` で最終集計を行っている。