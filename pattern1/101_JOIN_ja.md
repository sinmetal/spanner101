# JOIN

Pattern2は [インターリーブ](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child) していないスキーマ構成。
Albums TableはSingers TableはSingerIdで繋げることができるが、Spannerから見ると関連性はない。

## Sample Dataの追加

JOINするクエリのプロファイルを見る

```
cat ./dml/101_JOIN/sample_data.sql
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/101_JOIN/sample_data.sql
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
gcloud spanner cli $DB1 --instance=$CLOUDSDK_SPANNER_INSTANCE --project=$CLOUDSDK_CORE_PROJECT < ./dml/101_JOIN/query1.sql
+-----+--------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID  | Query_Execution_Plan                                                                       | Rows_Returned | Executions | Total_Latency |
+-----+--------------------------------------------------------------------------------------------+---------------+------------+---------------+
|   0 | Serialize Result                                                                           | 2             | 1          | 0.11 msecs    |
|   1 | +- Cross Apply                                                                             | 2             | 1          | 0.1 msecs     |
|  *2 |    +- [Input] Distributed Union (distribution_table: Singers, split_ranges_aligned: false) | 1             | 1          | 0.06 msecs    |
|   3 |    |  +- Local Distributed Union                                                           | 1             | 1          | 0.05 msecs    |
|  *4 |    |     +- Filter Scan (seekable_key_size: 1)                                             | 1             | 1          | 0.04 msecs    |
|   5 |    |        +- Table Scan (Table: Singers, scan_method: Scalar)                            | 1             | 1          | 0.04 msecs    |
| *16 |    +- [Map] Distributed Union (distribution_table: Albums, split_ranges_aligned: false)    | 2             | 1          | 0.04 msecs    |
|  17 |       +- Local Distributed Union                                                           | 2             | 1          | 0.03 msecs    |
| *18 |          +- Filter Scan                                                                    |               |            |               |
|  19 |             +- Table Scan (Table: Albums, scan_method: Scalar)                             | 2             | 1          | 0.03 msecs    |
+-----+--------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
  2: Split Range: ($SingerId = 1)
  4: Seek Condition: ($SingerId = 1)
 16: Split Range: ($SingerId_1 = 1)
 18: Seek Condition: ($SingerId_1 = 1)

2 rows in set (24.3 msecs)
timestamp:            2023-09-09T17:53:43.571124+09:00
cpu time:             4.9 msecs
rows scanned:         3 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230906_04_27_19UTC
```

`2` SingersとAlbums [Distributed Union](https://cloud.google.com/spanner/docs/query-execution-operators?hl=en#distributed-union) によって複数のSplitから `SingerId=1` のものだけ取り出され、その後、 [Cross Apply](https://cloud.google.com/spanner/docs/query-execution-operators?hl=en#cross-apply) によってJOINされる。
SingerId=1のSingersとAlbumsのRowは同じSplitにあるとは限らないので、Localでは完結しない。

# Refs

* https://spanner-hacks.apstn.dev/
* [Cloud Spanner でインターリーブテーブルを高速に取得する](https://medium.com/google-cloud-jp/cloud-spanner-%E3%81%A7%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%BC%E3%83%AA%E3%83%BC%E3%83%96%E3%83%86%E3%83%BC%E3%83%96%E3%83%AB%E3%82%92%E9%AB%98%E9%80%9F%E3%81%AB%E5%8F%96%E5%BE%97%E3%81%99%E3%82%8B-2a955b061d3)