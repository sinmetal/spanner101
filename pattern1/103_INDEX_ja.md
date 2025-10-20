# INDEX

Orders Tableから指定したUserのデータを取得する。

## Sample Dataの追加

Orders Tableに8行追加。
102_GROUPBYのサンプルデータを追加していない場合、先に102_GROUPBYのサンプルデータを追加してから行う。


```
cat ./dml/103_INDEX/sample_data.sql
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/sample_data.sql
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

欲しいのは5件だが、10件Scanしている。
今はTableに該当の行が10件しかないので、負荷が低いが、数が増えると不安がある。
以下のINDEXを作成しているが、利用されていない。

```
CREATE INDEX OrdersByUserIDAndCommitedAtDesc
    ON Orders (
        UserID,
        CommitedAt DESC
    );
```

INDEXを利用しなかった理由として考えられるのは以下

* `Amount` Columnがないので、Base TableとJOINする必要がある

## AmountをSTROINGしたINDEXを追加する

`Amount` ColumnをINDEX Tableに含めればBase TableとのJOINは必要なくなる。
`Amount` Column を [STORING](https://cloud.google.com/spanner/docs/secondary-indexes#storing-clause) してみよう。
STORINGするとINDEXのStorage Sizeが増えるので、必要に応じて追加しよう。
今回のケースではBase Tableと10件JOINするだけなのでさほど重いものではないが、このQueryが非常に高頻度で実行される場合は検討の余地がある。

STORINGのColumnの追加は既存のINDEXに対して行うことができるが、今回は新しいINDEXを作る。

``` create-index1.sql
CREATE INDEX OrdersByUserIDAndCommitedAtDescStoringAmount
ON Orders (
  UserID,
  CommitedAt DESC
) STORING (Amount);
```

```
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/create-index1.sql
```

## Orders Tableから指定したUserのデータを取得するクエリのプロファイルを再度見る

```
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/103_INDEX/query1.sql
```

```
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                                                 | Rows_Returned | Executions | Total_Latency |
+----+----------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Global Limit                                                                                                         | 5             | 1          | 0.07 msecs    |
| *1 | +- Distributed Union (distribution_table: OrdersByUserIDAndCommitedAtDescStoringAmount, split_ranges_aligned: false) | 5             | 1          | 0.07 msecs    |
|  2 |    +- Serialize Result                                                                                               | 5             | 1          | 0.06 msecs    |
|  3 |       +- Local Limit                                                                                                 | 5             | 1          | 0.05 msecs    |
|  4 |          +- Local Distributed Union                                                                                  | 5             | 1          | 0.05 msecs    |
| *5 |             +- Filter Scan                                                                                           |               |            |               |
|  6 |                +- Index Scan (Index: OrdersByUserIDAndCommitedAtDescStoringAmount, scan_method: Scalar)              | 5             | 1          | 0.05 msecs    |
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

Index Tableさえ参照すれば必要なColumnが揃うようになったので、Base TableとのJOINは必要なくなった。
