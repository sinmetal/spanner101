# JOIN

Pattern2は [インターリーブ](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child) しているスキーマ構成。
Albums TableはSingers Tableの子ども。

## Sample Dataの追加

Singers TableとAlbums Tableに1行ずつ追加

```
cat ./dml/101_JOIN/sample_data.sql
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/101_JOIN/sample_data.sql
```

## JOINするクエリのプロファイルを見る

Singers TableとAlbums TableのJOINを行うクエリのプロファイルを見る

``` query1.sql
EXPLAIN ANALYZE
SELECT * FROM Singers s
INNER JOIN Albums a ON s.SingerId = a.SingerId
WHERE s.SingerId = 1;
```

```
gcloud spanner cli $DB2 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/101_JOIN/query1.sql
+-----+-----------------------------------------------------------------------------+---------------+------------+---------------+
| ID  | Query_Execution_Plan                                                        | Rows_Returned | Executions | Total_Latency |
+-----+-----------------------------------------------------------------------------+---------------+------------+---------------+
|  *0 | Distributed Union (distribution_table: Singers, split_ranges_aligned: true) | 2             | 1          | 0.14 msecs    |
|   1 | +- Local Distributed Union                                                  | 2             | 1          | 0.12 msecs    |
|   2 |    +- Serialize Result                                                      | 2             | 1          | 0.11 msecs    |
|   3 |       +- Cross Apply                                                        | 2             | 1          | 0.1 msecs     |
|  *4 |          +- [Input] Filter Scan (seekable_key_size: 1)                      | 1             | 1          | 0.05 msecs    |
|   5 |          |  +- Table Scan (Table: Singers, scan_method: Scalar)             | 1             | 1          | 0.05 msecs    |
|  13 |          +- [Map] Local Distributed Union                                   | 2             | 1          | 0.05 msecs    |
| *14 |             +- Filter Scan                                                  |               |            |               |
|  15 |                +- Table Scan (Table: Albums, scan_method: Scalar)           | 2             | 1          | 0.05 msecs    |
+-----+-----------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
  0: Split Range: ($SingerId = 1)
  4: Seek Condition: ($SingerId = 1)
 14: Seek Condition: ($SingerId_1 = 1)

2 rows in set (6.29 msecs)
timestamp:            2023-09-09T17:42:56.704156+09:00
cpu time:             5.12 msecs
rows scanned:         3 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230906_07_18_51UTC
```

インターリーブしてないPattern1と比べてSingersとAlbumsをJOINする部分がLocalで完結するようになっています。
複数のマシンで実行されるものは [Distributed operators](https://cloud.google.com/spanner/docs/query-execution-operators?hl=en#distributed_operators) ですが、 [Cross Apply](https://cloud.google.com/spanner/docs/query-execution-operators?hl=en#cross-apply) の下にはありません。
インターリーブを組むことで、JOINを単一マシンで完結させることができています。

# Refs

* https://spanner-hacks.apstn.dev/
* [Cloud Spanner でインターリーブテーブルを高速に取得する](https://medium.com/google-cloud-jp/cloud-spanner-%E3%81%A7%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%BC%E3%83%AA%E3%83%BC%E3%83%96%E3%83%86%E3%83%BC%E3%83%96%E3%83%AB%E3%82%92%E9%AB%98%E9%80%9F%E3%81%AB%E5%8F%96%E5%BE%97%E3%81%99%E3%82%8B-2a955b061d3)