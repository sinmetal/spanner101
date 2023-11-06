# INDEX Pattern1

```
EXPLAIN ANALYZE
SELECT
  UserID,
  OrderID,
  Amount,
  CommitedAt
FROM Orders
WHERE UserID = "ruby"
ORDER BY CommitedAt DESC
LIMIT 10
```

```
+-----+---------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID  | Query_Execution_Plan                                                                                    | Rows_Returned | Executions | Total_Latency |
+-----+---------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|   0 | Global Limit                                                                                            | 10            | 1          | 3.89 msecs    |
|  *1 | +- Distributed Union (distribution_table: UserIDAndCommitedAtDescByOrders, split_ranges_aligned: false) | 10            | 1          | 3.89 msecs    |
|  *2 |    +- Distributed Cross Apply (order_preserving: true)                                                  | 10            | 1          | 3.87 msecs    |
|   3 |       +- [Input] Create Batch                                                                           |               |            |               |
|   4 |       |  +- Compute Struct                                                                              | 10            | 1          | 3.46 msecs    |
|   5 |       |     +- Local Limit                                                                              | 10            | 1          | 3.44 msecs    |
|   6 |       |        +- Local Distributed Union                                                               | 10            | 1          | 3.44 msecs    |
|  *7 |       |           +- Filter Scan                                                                        |               |            |               |
|   8 |       |              +- Index Scan (Index: UserIDAndCommitedAtDescByOrders, scan_method: Scalar)        | 10            | 1          | 3.43 msecs    |
|  21 |       +- [Map] Serialize Result                                                                         | 10            | 1          | 0.3 msecs     |
|  22 |          +- Cross Apply                                                                                 | 10            | 1          | 0.29 msecs    |
|  23 |             +- [Input] KeyRangeAccumulator                                                              |               |            |               |
|  24 |             |  +- Batch Scan (Batch: $v2, scan_method: Scalar)                                          |               |            |               |
|  28 |             +- [Map] Local Distributed Union                                                            | 10            | 10         | 0.26 msecs    |
| *29 |                +- Filter Scan (seekable_key_size: 1)                                                    | 10            | 10         | 0.25 msecs    |
|  30 |                   +- Table Scan (Table: Orders, scan_method: Scalar)                                    | 10            | 10         | 0.24 msecs    |
+-----+---------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
  1: Split Range: ($UserID = 'ruby')
  2: Split Range: ($OrderID' = $OrderID)
  7: Seek Condition: ($UserID = 'ruby')
 29: Seek Condition: ($OrderID' = $batched_OrderID)

10 rows in set (17.24 msecs)
timestamp:            2023-11-02T17:21:11.962571+09:00
cpu time:             9.19 msecs
rows scanned:         20 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20231028_13_47_28UTC
```

`UserIDAndCommitedAtDescByOrders` INDEXを参照して該当の10件を抽出後、 `Amount` Columnを取得するためにBase TableとJOINしている。

### STROINGしてみる

`Amount` ColumnをINDEX Tableに含めればBase TableとのJOINは必要なくなる。
`Amount` Column を [STORING](https://cloud.google.com/spanner/docs/secondary-indexes#storing-clause) してみよう。
STORINGするとINDEXのStorage Sizeが増えるので、必要に応じて追加しよう。
今回のケースではBase Tableと10件JOINするだけなのでさほど重いものではないが、このQueryが非常に高頻度で実行される場合は検討の余地がある。

```
CREATE INDEX UserIDAndCommitedAtDescStoringAmountByOrders
ON Orders (
  UserID,
  CommitedAt DESC
) STORING (Amount);
```

```
EXPLAIN ANALYZE
SELECT
  UserID,
  OrderID,
  Amount,
  CommitedAt
FROM Orders
WHERE UserID = "ruby"
ORDER BY CommitedAt DESC
LIMIT 10
```

```
spanner-cli -p gcpug-public-spanner -i merpay-sponsored-instance -d $DB1 -e "$(cat query.sql)" -t
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                                                 | Rows_Returned | Executions | Total_Latency |
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Global Limit                                                                                                         | 10            | 1          | 3.04 msecs    |
| *1 | +- Distributed Union (distribution_table: UserIDAndCommitedAtDescStoringAmountByOrders, split_ranges_aligned: false) | 10            | 1          | 3.03 msecs    |
|  2 |    +- Serialize Result                                                                                               | 10            | 1          | 3.02 msecs    |
|  3 |       +- Local Limit                                                                                                 | 10            | 1          | 3.01 msecs    |
|  4 |          +- Local Distributed Union                                                                                  | 10            | 1          | 3.01 msecs    |
| *5 |             +- Filter Scan                                                                                           |               |            |               |
|  6 |                +- Index Scan (Index: UserIDAndCommitedAtDescStoringAmountByOrders, scan_method: Scalar)              | 10            | 1          | 3 msecs       |
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 1: Split Range: ($UserID = 'ruby')
 5: Seek Condition: ($UserID = 'ruby')

10 rows in set (37.59 msecs)
timestamp:            2023-11-06T13:12:41.408123+09:00
cpu time:             4.57 msecs
rows scanned:         10 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20231104_13_30_53UTC
```

Index Tableさえ参照すれば必要なColumnが揃うようになったので、Base TableとのJOINは必要なくなった
