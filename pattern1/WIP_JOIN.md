# Pattern1

[インターリーブテーブルを高速に取得する](https://medium.com/google-cloud-jp/cloud-spanner-%E3%81%A7%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%BC%E3%83%AA%E3%83%BC%E3%83%96%E3%83%86%E3%83%BC%E3%83%96%E3%83%AB%E3%82%92%E9%AB%98%E9%80%9F%E3%81%AB%E5%8F%96%E5%BE%97%E3%81%99%E3%82%8B-2a955b061d3)


```
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("gold", "gold", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("silver", "silver", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("dia", "dia", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("ruby", "ruby", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("sapphire", "sapphire", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
```

```
INSERT INTO Orders(OrderID, UserID, Amount, CommitedAt) VALUES ("10ac9c3c-2e21-460e-be22-4527c11c1285","gold",100, PENDING_COMMIT_TIMESTAMP());
```

```
EXPLAIN ANALYZE
WITH
  TargetOrders AS (
  SELECT
    Orders.OrderID,
    Orders.CommitedAt,
  FROM
    Orders
  WHERE
    Orders.UserID = "ruby"
  ORDER BY
    Orders.CommitedAt DESC
  LIMIT
    30 )
SELECT
  Orders.OrderID,
  Orders.CommitedAt,
   ARRAY(
     SELECT STRUCT<OrderDetailID STRING, ItemID STRING, Price INT64, Quantity INT64>
     (OrderDetailID,
      ItemID,
      Price,
      Quantity)) AS OrderDetails
FROM TargetOrders AS Orders JOIN OrderDetails ON Orders.OrderID = OrderDetails.OrderID
```

```
spanner-cli -p gcpug-public-spanner -i merpay-sponsored-instance -d sinmetal -e "$(cat query.sql)" -t                                                                                                                                                                1 ↵
```

```
+-----+------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID  | Query_Execution_Plan                                                                                             | Rows_Returned | Executions | Total_Latency |
+-----+------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|  *0 | Distributed Cross Apply                                                                                          | 469           | 1          | 32.93 msecs   |
|   1 | +- [Input] Create Batch                                                                                          |               |            |               |
|   2 | |  +- Compute Struct                                                                                             | 30            | 1          | 3 msecs       |
|   3 | |     +- Global Limit                                                                                            | 30            | 1          | 2.97 msecs    |
|  *4 | |        +- Distributed Union (distribution_table: UserIDAndCommitedAtDescByOrders, split_ranges_aligned: false) | 30            | 1          | 2.97 msecs    |
|   5 | |           +- Local Limit                                                                                       | 30            | 1          | 2.95 msecs    |
|   6 | |              +- Local Distributed Union                                                                        | 30            | 1          | 2.95 msecs    |
|  *7 | |                 +- Filter Scan                                                                                 |               |            |               |
|   8 | |                    +- Index Scan (Index: UserIDAndCommitedAtDescByOrders, scan_method: Scalar)                 | 30            | 1          | 2.94 msecs    |
|  24 | +- [Map] Serialize Result                                                                                        | 469           | 1          | 29.71 msecs   |
|  25 |    +- Compute Struct                                                                                             | 469           | 1          | 29.24 msecs   |
|  26 |       +- Cross Apply                                                                                             | 469           | 1          | 29.05 msecs   |
|  27 |          +- [Input] KeyRangeAccumulator                                                                          |               |            |               |
|  28 |          |  +- Batch Scan (Batch: $v3, scan_method: Scalar)                                                      |               |            |               |
|  31 |          +- [Map] Local Distributed Union                                                                        | 469           | 30         | 28.98 msecs   |
| *32 |             +- Filter Scan                                                                                       |               |            |               |
|  33 |                +- Table Scan (Table: OrderDetails, scan_method: Scalar)                                          | 469           | 30         | 28.93 msecs   |
+-----+------------------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
  0: Split Range: ($OrderID_3 = $OrderID)
  4: Split Range: ($UserID = 'ruby')
  7: Seek Condition: ($UserID = 'ruby')
 32: Seek Condition: ($OrderID_3 = $batched_OrderID)

469 rows in set (60.15 msecs)
timestamp:            2023-09-07T20:12:35.261622+09:00
cpu time:             11.31 msecs
rows scanned:         499 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230906_04_27_19UTC
```