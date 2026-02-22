##### Build the db-cleaner
```shell
docker build -f cmd/db-cleaner/Dockerfile -t data-platform/db-cleaner:0.0.1 .
```
##### Run the db-cleaner
```shell
docker run --name db-cleaner \
 -e DATA_CATALOG_URL="https://openmetadata.example.com/api" \
 -e DATA_CATALOG_TOKEN="eyJraWQiOiJHYj" \
 -e DATA_CATALOG_DATABASE_SCHEMA="service.database.schema" \
 -e DATABASE_URL="postgres://username:password@postgresql.example.com:5432/database?sslmode=disable" \
 -e DATABASE_SCHEMA="schema" \
 data-platform/db-cleaner:0.0.1
```
