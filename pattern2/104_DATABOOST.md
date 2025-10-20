# DataBoost

https://cloud.google.com/spanner/docs/databoost/databoost-overview

Execute [Federated Query](https://cloud.google.com/bigquery/docs/cloud-spanner-federated-queries) using DataBoost from BigQuery.

```
properties2=$(echo '{"database":"projects/'$CLOUDSDK_CORE_PROJECT'/instances/'$CLOUDSDK_SPANNER_INSTANCE'/databases/'$DB2'", "useParallelism":true, "useDataBoost": true}')
connection_name2=$(printf "spanner_%s_%s" $CLOUDSDK_SPANNER_INSTANCE $DB2)
bq mk --project_id $CLOUDSDK_CORE_PROJECT --connection --connection_type='CLOUD_SPANNER' --location='us-central1' \
--properties=$properties $connection_name2
```

# JOIN

Since Orders are interleaved as children of Users, they can be completed in the same Partition when joining.

```
bq query --use_legacy_sql=false << EOS
SELECT * FROM EXTERNAL_QUERY(
  '$CLOUDSDK_CORE_PROJECT.us-central1.$connection_name2',
  'SELECT Users.UserID,Orders.OrderID FROM Users JOIN Orders ON Users.UserID = Orders.UserID') AS UserOrders
EOS
```

# Memo

miscellaneous notes from sinmetal

If it is an INDEX that is interleaved with Users, it can be executed even if it is referenced by DataBoost.
However, since a query like the one below results in a Residual Condition, it doesn't contribute much to lowering the processing load, so wouldn't it be better to just do a Table Full Scan? Maybe it feels like that?

```
SELECT * FROM EXTERNAL_QUERY(
  'gcpug-public-spanner.us-central1.spanner_sinmetal2',
  '''SELECT
       Users.UserID,
       Orders.OrderID,
       Orders.Amount,
       Orders.CommitedAt
     FROM Users JOIN Orders@{FORCE_INDEX=OrdersByUserIDAndCommitedAtDescParentUsers} ON Users.UserID = Orders.UserID
     WHERE FORMAT_TIMESTAMP("%Y%m",Orders.CommitedAt, "Asia/Tokyo") = "202309"''') AS UserOrders
```

```
+-----+---------------------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID  | Query_Execution_Plan                                                                                                            | Rows_Returned | Executions | Total_Latency |
+-----+---------------------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|   0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)                                                       | 19573         | 1          | 1.01 secs     |
|   1 | +- Local Distributed Union                                                                                                      | 19573         | 1          | 1.01 secs     |
|   2 |    +- Serialize Result                                                                                                          | 19573         | 1          | 1.01 secs     |
|   3 |       +- Cross Apply                                                                                                            | 19573         | 1          | 994.74 msecs  |
|  *4 |          +- [Input] Filter Scan                                                                                                 |               |            |               |
|   5 |          |  +- Index Scan (Full scan: true, Index: OrdersByUserIDAndCommitedAtDescStoredAmountParentUsers, scan_method: Scalar) | 19573         | 1          | 34.63 msecs   |
|  17 |          +- [Map] Local Distributed Union                                                                                       | 19573         | 19573      | 952.32 msecs  |
| *18 |             +- Filter Scan (seekable_key_size: 1)                                                                               | 19573         | 19573      | 941.31 msecs  |
|  19 |                +- Table Scan (Table: Users, scan_method: Scalar)                                                                | 19573         | 19573      | 924.97 msecs  |
+-----+---------------------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
  4: Residual Condition: (FORMAT_TIMESTAMP('%Y%m', $CommitedAt, 'Asia/Tokyo') = '202309')
 18: Seek Condition: ($UserID = $UserID_1)

19573 rows in set (1.71 secs)
timestamp:            2023-09-11T17:27:30.192549+09:00
cpu time:             349.2 msecs
rows scanned:         39146 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230906_07_18_51UTC
```
