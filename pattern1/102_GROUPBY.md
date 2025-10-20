# GROUP BY

Aggregate the amount of Orders Table for each User.
The PK of the Orders Table is OrderID.

## Sample Dataの追加

Add 5 rows to Users Table and 6 rows to Orders Table

```
cat ./dml/102_GROUPBY/sample_data.sql
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/102_GROUPBY/sample_data.sql
```

## View the profile of a query that aggregates by GROUP BY for each user

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
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/102_GROUPBY/query1.sql
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

`Hash Aggregate` is introduced for GROUP BY.
`Local Hash Aggregate` -> `Distributed Union` -> `Global Hash Aggregate`, so after aggregating on each machine, they are collected and aggregated again.
In this example, there are only 5 types of UserID, so it can be completed quickly, but if there are many users, it will be difficult to create a Hash Table and tally it.
