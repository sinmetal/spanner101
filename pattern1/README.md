# Pattern1

Pattern1は [インターリーブ](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child) していないスキーマ構成。

``` Pattern1用のDBを作成する
export CLOUDSDK_CORE_PROJECT=gcpug-public-spanner
export CLOUDSDK_SPANNER_INSTANCE=spanner101
export DB1=sample1
gcloud spanner databases create $DB1 --ddl "$(cat ./ddl/ddl.sql)"
```

DB1の名前は他の人と重複しないように別のものにしておくとよいです。