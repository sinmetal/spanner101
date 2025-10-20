# Spanner101

Spanner hands-on materials

[Interleaving](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child) is an important and unique feature when designing Spanner tables. Hands-on to touch the center.

* pattern1 no interleaving
* pattern2 with interleaving
* pattern3 bonus

Since the execution plan is viewed using [https://cloud.google.com/spanner/docs/spanner-cli?hl=en](https://cloud.google.com/spanner/docs/spanner-cli?hl=en), please install the gcloud cli first if it is not already installed on your machine.

## Singers

Experience interleaving with a simple query using Singers and Albums, which are also available in the official documentation.

## DataBoost

Hands-on on [DataBoost](https://cloud.google.com/spanner/docs/databoost/databoost-overview) to run queries without putting any load on Spanner Instance.

The queries that can be executed with DataBoost are limited to [PartitionQuery](https://cloud.google.com/spanner/docs/reads?hl=en#read_data_in_parallel).
We will experience the changes in the query that can be made into PartitionQuery with and without interleaving.