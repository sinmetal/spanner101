# GROUP BY

Orders TableのAmountをUserごとに集計する。
Orders TableのPKはOrderIDになっている。

## Sample Dataの追加

Users Tableに5行、Orders Tableに6行を追加

```
cat ./dml/102_GROUPBY/sample_data.sql
spanner-cli -p $CLOUDSDK_CORE_PROJECT -i $CLOUDSDK_SPANNER_INSTANCE -d $DB1 -e "$(cat ./dml/102_GROUPBY/sample_data.sql)" -t
```

## UserごとにGROUP BYで集計を行うクエリのプロファイルを見る

``` query1.sql
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
spanner-cli -p $CLOUDSDK_CORE_PROJECT -i $CLOUDSDK_SPANNER_INSTANCE -d $DB1 -e "$(cat ./dml/102_GROUPBY/query1.sql)" -t
```

```
+----+------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                               | Rows_Returned | Executions | Total_Latency |
+----+------------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Serialize Result                                                                   | 5             | 1          | 0.09 msecs    |
|  1 | +- Global Hash Aggregate                                                           | 5             | 1          | 0.09 msecs    |
|  2 |    +- Distributed Union (distribution_table: Orders, split_ranges_aligned: false)  | 5             | 1          | 0.08 msecs    |
|  3 |       +- Local Hash Aggregate                                                      | 5             | 1          | 0.07 msecs    |
|  4 |          +- Local Distributed Union                                                | 6             | 1          | 0.06 msecs    |
|  5 |             +- Table Scan (Full scan: true, Table: Orders, scan_method: Automatic) | 6             | 1          | 0.05 msecs    |
+----+------------------------------------------------------------------------------------+---------------+------------+---------------+
5 rows in set (3.08 msecs)
timestamp:            2024-02-27T15:05:09.366162+09:00
cpu time:             3.05 msecs
rows scanned:         6 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20240226_09_13_34UTC
```

GROUP BYのために `Hash Aggregate` が登場。
`Local Hash Aggregate` -> `Distributed Union` -> `Global Hash Aggregate` となっているので、各マシンで集計後、それらを集めて、再度集計している。
この例だとUserIDの種類が5つしかないので、すぐ終わっているが、Userが多いとHash Table作って集計するのが大変になる。
