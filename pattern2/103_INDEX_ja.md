# INDEX

Orders Tableから指定したUserのデータを取得する。

## Sample Dataの追加

Orders Tableに8行追加。
102_GROUPBYのサンプルデータを追加していない場合、先に102_GROUPBYのサンプルデータを追加してから行う。

```
cat ./dml/103_INDEX/sample_data.sql
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/sample_data.sql
```

## Orders Tableから指定したUserのデータを取得するクエリのプロファイルを見る

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

欲しいのは5件だが、10件Scanしている。
今はTableに該当の行が10件しかないので、負荷が低いが、数が増えると不安がある。
pattern1と比べると `Sort Limit` がLocal Distributed Unionのすぐ後にあるので、1つのマシンでSort, Limitが完結できている。
以下のINDEXを作成しているが、利用されていない。

```
CREATE INDEX UserIDAndCommitedAtDescParentUsersByOrders
    ON Orders (
        UserID,
        CommitedAt DESC
    ), INTERLEAVE IN Users;
```

INDEXを利用しなかった理由として考えられるのは以下

* `Amount` Columnがないので、Base TableとJOINする必要がある

## AmountをSTROINGしたINDEXを追加する

`Amount` ColumnをINDEX Tableに含めればBase TableとのJOINは必要なくなる。
`Amount` Column を [STORING](https://cloud.google.com/spanner/docs/secondary-indexes#storing-clause) してみよう。
STORINGするとINDEXのStorage Sizeが増えるので、必要に応じて追加しよう。
今回のケースは件数が多くなると負荷が上がっていく可能性があるので、検討の余地がある。

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

## Orders Tableから指定したUserのデータを取得するクエリのプロファイルを再度見る

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

Index Tableさえ参照すれば必要なColumnが揃うようになったので、Base TableとのJOINは必要なくなった。
pattern1と比べると `Global Limit` が `Local Distributed Union` のすぐ後に来ているので、1つのマシンでLimitが完結できている。