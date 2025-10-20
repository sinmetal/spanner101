# Pattern2

Pattern2は [インターリーブ](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child) しているスキーマ構成。

``` Pattern1用のDBを作成する
export CLOUDSDK_CORE_PROJECT=gcpug-public-spanner
export CLOUDSDK_SPANNER_INSTANCE=spanner101
export DB2=sample2
gcloud spanner databases create $DB2 --ddl-file=./ddl/ddl.sql
```

DB2の名前は他の人と重複しないように別のものにしておくとよいです。