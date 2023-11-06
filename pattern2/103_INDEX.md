# INDEX Pattern2

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
export DB1={YOUR DB}
spanner-cli -p gcpug-public-spanner -i merpay-sponsored-instance -d $DB2 -e "$(cat query.sql)" -t
+-----+-------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID  | Query_Execution_Plan                                                                                  | Rows_Returned | Executions | Total_Latency |
+-----+-------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|  *0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)                             | 10            | 1          | 2.85 msecs    |
|   1 | +- Serialize Result                                                                                   | 10            | 1          | 2.84 msecs    |
|   2 |    +- Global Limit                                                                                    | 10            | 1          | 2.82 msecs    |
|   3 |       +- Local Distributed Union                                                                      | 10            | 1          | 2.82 msecs    |
|   4 |          +- Cross Apply                                                                               | 10            | 1          | 2.81 msecs    |
|  *5 |             +- [Input] Filter Scan                                                                    |               |            |               |
|   6 |             |  +- Index Scan (Index: UserIDAndCommitedAtDescByOrdersParentUsers, scan_method: Scalar) | 10            | 1          | 0.34 msecs    |
|  14 |             +- [Map] Local Distributed Union                                                          | 10            | 10         | 2.45 msecs    |
| *15 |                +- Filter Scan (seekable_key_size: 2)                                                  | 10            | 10         | 2.44 msecs    |
|  16 |                   +- Table Scan (Table: Orders, scan_method: Scalar)                                  | 10            | 10         | 2.43 msecs    |
+-----+-------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
  0: Split Range: ($UserID = 'ruby')
  5: Seek Condition: ($UserID = 'ruby')
 15: Seek Condition: ($UserID' = $UserID) AND ($OrderID' = $OrderID)

10 rows in set (11.46 msecs)
timestamp:            2023-11-06T14:13:52.786129+09:00
cpu time:             8.67 msecs
rows scanned:         20 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20231106_02_31_17UTC
```

`UserIDAndCommitedAtDescByOrdersParentUsers` がちょうどよいINDEXなので、それを参照し、INDEXにないColumnを取得するためにBaseTableとJOINしている。
`UserIDAndCommitedAtDescByOrdersParentUsers` は `Users` Tableの子どもなので、JOINは `Local Distributed Union` で完結している

### STROINGしてみる

`Amount` ColumnをINDEX Tableに含めればBase TableとのJOINは必要なくなる。
`Amount` Column を [STORING](https://cloud.google.com/spanner/docs/secondary-indexes#storing-clause) してみよう。
STORINGするとINDEXのStorage Sizeが増えるので、必要に応じて追加しよう。
今回のケースではBase Tableと10件JOINするだけなのでさほど重いものではないが、このQueryが非常に高頻度で実行される場合は検討の余地がある。

```
CREATE INDEX UserIDAndCommitedAtDescStoringAmountParentUsersByOrders
ON Orders (
  UserID,
  CommitedAt DESC
) STORING (Amount), INTERLEAVE IN Users;
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
spanner-cli -p gcpug-public-spanner -i merpay-sponsored-instance -d $DB2 -e "$(cat query.sql)" -t
+----+-----------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                                            | Rows_Returned | Executions | Total_Latency |
+----+-----------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| *0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)                                       | 10            | 1          | 0.11 msecs    |
|  1 | +- Serialize Result                                                                                             | 10            | 1          | 0.1 msecs     |
|  2 |    +- Global Limit                                                                                              | 10            | 1          | 0.09 msecs    |
|  3 |       +- Local Distributed Union                                                                                | 10            | 1          | 0.09 msecs    |
| *4 |          +- Filter Scan                                                                                         |               |            |               |
|  5 |             +- Index Scan (Index: UserIDAndCommitedAtDescStoringAmountParentUsersByOrders, scan_method: Scalar) | 10            | 1          | 0.08 msecs    |
+----+-----------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 0: Split Range: ($UserID = 'ruby')
 4: Seek Condition: ($UserID = 'ruby')

10 rows in set (5.48 msecs)
timestamp:            2023-11-06T14:29:20.905981+09:00
cpu time:             5.42 msecs
rows scanned:         10 rows
deleted rows scanned: 0 rows
optimizer version:    6
optimizer statistics: auto_20231106_02_31_17UTC
```