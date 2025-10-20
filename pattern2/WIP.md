```
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("gold", "gold", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("silver", "silver", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("dia", "dia", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("ruby", "ruby", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("sapphire", "sapphire", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders(UserID, OrderID, Amount, CommitedAt) VALUES ("gold","10ac9c3c-2e21-460e-be22-4527c11c1285", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO OrderDetails(UserID, OrderID, OrderDetailID, ItemID, Price, Quantity, CommitedAt) VALUES("gold", "10ac9c3c-2e21-460e-be22-4527c11c1285", 1, "pen", 100, 1, PENDING_COMMIT_TIMESTAMP());
```

```
SELECT
  Orders.UserID,
  Orders.OrderID,
  Orders.Amount,
  Orders.CommitedAt,
  OrderDetails.OrderDetailID,
  OrderDetails.ItemID,
  OrderDetails.Price,
  OrderDetails.Quantity
FROM
  Orders JOIN OrderDetails ON Orders.OrderID = OrderDetails.OrderID
WHERE
  Orders.CommitedAt >= TIMESTAMP('2023-09-01 00:00:00', "Asia/Tokyo")
  AND Orders.CommitedAt < TIMESTAMP('2023-10-01 00:00:00', "Asia/Tokyo")
```


```
EXPLAIN ANALYZE
WITH
  TargetOrders AS (
  SELECT
    Orders.UserID,
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
FROM TargetOrders AS Orders JOIN OrderDetails ON Orders.UserID = OrderDetails.UserID AND Orders.OrderID = OrderDetails.OrderID
```

```
gcloud spanner cli sinmetal2 --instance=merpay-sponsored-instance --project=gcpug-public-spanner < query.sql
```

```
+-----+-------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID  | Query_Execution_Plan                                                                                  | Rows_Returned | Executions | Total_Latency |
+-----+-------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
|  *0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)                             | 527           | 1          | 257.53 msecs  |
|   1 | +- Serialize Result                                                                                   | 527           | 1          | 257.46 msecs  |
|   2 |    +- Cross Apply                                                                                     | 527           | 1          | 256.9 msecs   |
|   3 |       +- [Input] Global Limit                                                                         | 30            | 1          | 254.81 msecs  |
|   4 |       |  +- Local Distributed Union                                                                   | 30            | 1          | 254.81 msecs  |
|  *5 |       |     +- Filter Scan                                                                            |               |            |               |
|   6 |       |        +- Index Scan (Index: UserIDAndCommitedAtDescByOrdersParentUsers, scan_method: Scalar) | 30            | 1          | 254.8 msecs   |
|  15 |       +- [Map] Local Distributed Union                                                                | 527           | 30         | 2.04 msecs    |
|  16 |          +- Compute Struct                                                                            | 527           | 30         | 1.97 msecs    |
| *17 |             +- Filter Scan                                                                            |               |            |               |
|  18 |                +- Table Scan (Table: OrderDetails, scan_method: Scalar)                               | 527           | 30         | 1.75 msecs    |
+-----+-------------------------------------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
  0: Split Range: ($UserID = 'ruby')
  5: Seek Condition: ($UserID = 'ruby')
 17: Seek Condition: (($UserID_3 = 'ruby') AND ($OrderID_3 = $OrderID))

527 rows in set (309.14 msecs)
timestamp:            2023-09-07T20:21:24.925704+09:00
cpu time:             10.08 msecs
rows scanned:         557 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230906_07_18_51UTC
```

```

```