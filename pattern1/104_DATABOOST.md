# DataBoost Pattern1

https://cloud.google.com/spanner/docs/databoost/databoost-overview

BigQueryからDataBoostを利用して [Federated Query](https://cloud.google.com/bigquery/docs/cloud-spanner-federated-queries) を実行する。
SpannerとBigQueryは同じRegionである必要がある点に注意。

```
bq mk --project_id gcpug-public-spanner --connection --connection_type='CLOUD_SPANNER' --location='us-central1' \
--properties='{"database":"projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal1", "useParallelism":true, "useDataBoost": true}' spanner_sinmetal1
```

# JOIN

DataBoostでSpannerに対して実行するQueryはPartitionQueryとして実行される。
そのため、実行できるQueryに制限が存在する。
pattern1の場合、UsersとOrdersはインターリーブされてないので、JOINしようとすると別のPartitionにRowが跨ってしまうため、実行できない

```
# このクエリは実行できない
SELECT * FROM EXTERNAL_QUERY(
  'gcpug-public-spanner.us-central1.spanner_sinmetal1',
  'SELECT Users.UserID,Orders.OrderID FROM Users JOIN Orders ON Users.UserID = Orders.UserID') AS UserOrders
```