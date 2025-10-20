# Pattern1

Pattern1 is a schema configuration that is not [interleaved](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child).

```Create a DB for Pattern1
export CLOUDSDK_CORE_PROJECT=gcpug-public-spanner
export CLOUDSDK_SPANNER_INSTANCE=spanner101
export DB1=sample1
gcloud spanner databases create $DB1 --ddl-file=./ddl/ddl.sql
````

It is a good idea to use a different name for DB1 to avoid duplication with others.