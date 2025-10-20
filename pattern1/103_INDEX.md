# INDEX

Get the data of the specified User from the Orders Table.

## Add sample data

Added 8 rows to Orders Table.
If you have not added sample data for 102_GROUPBY, add sample data for 102_GROUPBY first.

```
cat ./dml/103_INDEX/sample_data.sql
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/sample_data.sql
```

## View the profile of the query that retrieves data for the specified User from the Orders Table

gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/query1.sql
```

```
+----+---------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                                          | Rows_Returned | Executions | Total_Latency |
+----+---------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Global Limit                                                                                                  | 5             | 1          | 0.09 msecs    |
| *1 | +- Distributed Union (distribution_table: Orders, preserve_subquery_order: true, split_ranges_aligned: false) | 5             | 1          | 0.09 msecs    |
|  2 |    +- Serialize Result                                                                                        | 5             | 1          | 0.07 msecs    |
|  3 |       +- Local Sort Limit                                                                                     | 5             | 1          | 0.07 msecs    |
|  4 |          +- Local Distributed Union                                                                           | 10            | 1          | 0.06 msecs    |
| *5 |             +- Filter Scan                                                                                    |               |            |               |
|  6 |                +- Table Scan (Table: Orders, scan_method: Automatic)                                          | 10            | 1          | 0.05 msecs    |
+----+---------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 1: Split Range: ($UserID = 'ruby')
 5: Seek Condition: ($UserID = 'ruby')

5 rows in set (12.75 msecs)
timestamp:            2024-02-27T20:09:21.648723+09:00
cpu time:             11.71 msecs
rows scanned:         10 rows
deleted rows scanned: 0 rows
optimizer version:    6
```

I want 5 items, but I have scanned 10 items.
Currently, there are only 10 relevant rows in the table, so the load is low, but I'm worried if the number increases.
The following INDEX has been created, but it is not used.

```
CREATE INDEX UserIDAndCommitedAtDescByOrders
    ON Orders (
        UserID,
        CommitedAt DESC
    );
```

Possible reasons for not using INDEX are as follows:

* There is no `Amount` Column, so you need to JOIN with Base Table

## Add INDEX with STROING Amount

If you include the `Amount` Column in the INDEX Table, there is no need to JOIN with the Base Table.
Let's [STORING](https://cloud.google.com/spanner/docs/secondary-indexes#storing-clause) the `Amount` Column.
STORING will increase the Storage Size of INDEX, so add it as necessary.
In this case, it's not that heavy since it's just a 10-item JOIN with the Base Table, but there's room for consideration if this query is executed very frequently.

Adding a STORING Column can be done to an existing INDEX, but this time we will create a new INDEX.

``` create-index1.sql
CREATE INDEX UserIDAndCommitedAtDescStoringAmountByOrders
ON Orders (
  UserID,
  CommitedAt DESC
) STORING (Amount);
```

```
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/create-index1.sql
```

## Look again at the profile of the query that retrieves data for the specified User from the Orders Table

```
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/query1.sql
```

```
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                                                 | Rows_Returned | Executions | Total_Latency |
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Global Limit                                                                                                         | 5             | 1          | 0.07 msecs    |
| *1 | +- Distributed Union (distribution_table: UserIDAndCommitedAtDescStoringAmountByOrders, split_ranges_aligned: false) | 5             | 1          | 0.07 msecs    |
|  2 |    +- Serialize Result                                                                                               | 5             | 1          | 0.06 msecs    |
|  3 |       +- Local Limit                                                                                                 | 5             | 1          | 0.05 msecs    |
|  4 |          +- Local Distributed Union                                                                                  | 5             | 1          | 0.05 msecs    |
| *5 |             +- Filter Scan                                                                                           |               |            |               |
|  6 |                +- Index Scan (Index: UserIDAndCommitedAtDescStoringAmountByOrders, scan_method: Scalar)              | 5             | 1          | 0.05 msecs    |
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 1: Split Range: ($UserID = 'ruby')
 5: Seek Condition: ($UserID = 'ruby')

5 rows in set (10.68 msecs)
timestamp:            2024-02-27T20:12:24.929053+09:00
cpu time:             10.64 msecs
rows scanned:         5 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20240227_11_07_41UTC
```

Since the necessary columns are now available just by referencing the Index Table, JOIN with the Base Table is no longer necessary.
