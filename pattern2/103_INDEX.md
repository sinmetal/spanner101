#INDEX

Get the data of the specified User from the Orders Table.

## Add Sample Data

Added 8 rows to Orders Table.
If you have not added sample data for 102_GROUPBY, add sample data for 102_GROUPBY first.

````
cat ./dml/103_INDEX/sample_data.sql
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/sample_data.sql
````

## View the profile of the query that retrieves data for the specified User from the Orders Table

``` query1.sql
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
```

```
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/query1.sql
```

```
+----+---------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                      | Rows_Returned | Executions | Total_Latency |
+----+---------------------------------------------------------------------------+---------------+------------+---------------+
| *0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true) | 5             | 1          | 0.1 msecs     |
|  1 | +- Serialize Result                                                       | 5             | 1          | 0.08 msecs    |
|  2 |    +- Global Sort Limit                                                   | 5             | 1          | 0.08 msecs    |
|  3 |       +- Local Distributed Union                                          | 10            | 1          | 0.06 msecs    |
| *4 |          +- Filter Scan                                                   |               |            |               |
|  5 |             +- Table Scan (Table: Orders, scan_method: Automatic)         | 10            | 1          | 0.06 msecs    |
+----+---------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 0: Split Range: ($UserID = 'ruby')
 4: Seek Condition: ($UserID = 'ruby')

5 rows in set (5.91 msecs)
timestamp:            2024-02-27T19:47:21.58119+09:00
cpu time:             4.94 msecs
rows scanned:         10 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20240227_05_47_04UTC
```

I want 5 items, but I have scanned 10 items.
Currently, there are only 10 relevant rows in the table, so the load is low, but I'm worried if the number increases.
Compared to pattern 1, `Sort Limit` is located immediately after Local Distributed Union, so Sort and Limit can be completed on one machine.
The following INDEX has been created, but it is not used.

```
CREATE INDEX UserIDAndCommitedAtDescParentUsersByOrders
    ON Orders (
        UserID,
        CommitedAt DESC
    ), INTERLEAVE IN Users;
```

Possible reasons for not using INDEX are as follows:

* There is no `Amount` Column, so you need to JOIN with Base Table

## Add INDEX with STROING Amount

If you include the `Amount` Column in the INDEX Table, there is no need to JOIN with the Base Table.
Let's [STORING](https://cloud.google.com/spanner/docs/secondary-indexes#storing-clause) the `Amount` Column.
STORING will increase the Storage Size of INDEX, so add it as necessary.
In this case, there is a possibility that the load will increase as the number of cases increases, so there is room for consideration.

``` create-index1.sql
CREATE INDEX UserIDAndCommitedAtDescStoringAmountParentUsersByOrders
ON Orders (
  UserID,
  CommitedAt DESC
) STORING (Amount), INTERLEAVE IN Users;
```

```
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/create-index1.sql
```

## Look again at the profile of the query that retrieves data for the specified User from the Orders Table

```
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/query1.sql
```

```
+----+-----------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                                            | Rows_Returned | Executions | Total_Latency |
+----+-----------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| *0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)                                       | 5             | 1          | 0.08 msecs    |
|  1 | +- Serialize Result                                                                                             | 5             | 1          | 0.07 msecs    |
|  2 |    +- Global Limit                                                                                              | 5             | 1          | 0.06 msecs    |
|  3 |       +- Local Distributed Union                                                                                | 5             | 1          | 0.06 msecs    |
| *4 |          +- Filter Scan                                                                                         |               |            |               |
|  5 |             +- Index Scan (Index: UserIDAndCommitedAtDescStoringAmountParentUsersByOrders, scan_method: Scalar) | 5             | 1          | 0.06 msecs    |
+----+-----------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 0: Split Range: ($UserID = 'ruby')
 4: Seek Condition: ($UserID = 'ruby')

5 rows in set (11.84 msecs)
timestamp:            2024-02-27T20:02:10.113711+09:00
cpu time:             10.84 msecs
rows scanned:         5 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20240227_05_47_04UTC
```

Since the necessary columns are now available just by referencing the Index Table, JOIN with the Base Table is no longer necessary.
Compared to pattern 1, `Global Limit` comes immediately after `Local Distributed Union`, so Limit can be completed on one machine.