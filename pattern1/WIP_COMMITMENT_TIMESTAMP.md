# Commitment Timestamp

SpannerにはTIMESTAMP型のColumnに [CommitのTimestampを利用できるようにするオプション](https://cloud.google.com/spanner/docs/commit-timestamp) を指定できる。
このColumnを利用して [期間のFilterした場合、I/Oの最適化が入る](https://cloud.google.com/spanner/docs/commit-timestamp?hl=en#optimize)。
INDEXを作らなくても、直近のデータを取得できるので便利。

```
EXPLAIN ANALYZE
SELECT
  *
FROM
  Orders
WHERE
  CommitedAt >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR)
```

```
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                         | Rows_Returned | Executions | Total_Latency |
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Distributed Union (distribution_table: Orders, split_ranges_aligned: false)  | 60            | 1          | 0.77 msecs    |
|  1 | +- Local Distributed Union                                                   | 60            | 1          | 0.74 msecs    |
|  2 |    +- Serialize Result                                                       | 60            | 1          | 0.73 msecs    |
| *3 |       +- Filter Scan                                                         |               |            |               |
|  4 |          +- Table Scan (Full scan: true, Table: Orders, scan_method: Scalar) | 60            | 1          | 0.7 msecs     |
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 3: Residual Condition: ($CommitedAt >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), 1, HOUR))

60 rows in set (4.98 msecs)
timestamp:            2023-09-13T17:50:11.939791+09:00
cpu time:             3.83 msecs
rows scanned:         781 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230909_06_42_57UTC
```