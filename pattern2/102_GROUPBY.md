# GROUP BY

Aggregate the amount of Orders Table for each User.
Orders Table is a child of Users Table, so PK is UserID, OrderID.

## Add Sample Data

Add 5 rows to Users Table and 6 rows to Orders Table

````
cat ./dml/102_GROUPBY/sample_data.sql
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/102_GROUPBY/sample_data.sql
````

## View the profile of a query that aggregates with GROUP BY for each User

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
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/102_GROUPBY/query1.sql
```

```
+----+---------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                            | Rows_Returned | Executions | Total_Latency |
+----+---------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)       | 5             | 1          | 0.1 msecs     |
|  1 | +- Serialize Result                                                             | 5             | 1          | 0.08 msecs    |
|  2 |    +- Stream Aggregate                                                          | 5             | 1          | 0.07 msecs    |
|  3 |       +- Local Distributed Union                                                | 6             | 1          | 0.07 msecs    |
|  4 |          +- Table Scan (Full scan: true, Table: Orders, scan_method: Automatic) | 6             | 1          | 0.06 msecs    |
+----+---------------------------------------------------------------------------------+---------------+------------+---------------+
5 rows in set (4.82 msecs)
timestamp:            2024-02-27T16:08:05.452253+09:00
cpu time:             3.86 msecs
rows scanned:         6 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20240227_05_47_04UTC
```

Since the UserID is at the beginning of the PK, the Orders are grouped together for each User.
Therefore, sequential aggregation is possible with Stream Aggregate.
Orders Table is a child of Users Table, so Orders for the same User are on the same machine.
Therefore, it is completed locally.