# Pattern2

Pattern2 is a schema structure that is [interleaved](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child).

```Create a DB for Pattern1
export CLOUDSDK_CORE_PROJECT=gcpug-public-spanner
export CLOUDSDK_SPANNER_INSTANCE=spanner101
export DB2=sample2
gcloud spanner databases create $DB2 --ddl "$(cat ./ddl/ddl.sql)"
````

It is a good idea to use a different name for DB2 to avoid duplication with others.