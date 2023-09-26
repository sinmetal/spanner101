# Pattern3

Pattern2から、Orders TableのPKを日時の降順に変更したもの。
指定したユーザのOrdersを降順に表示することが多い時にPK順に取得できる。

```
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("gold", "gold", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("silver", "silver", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("dia", "dia", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("ruby", "ruby", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("sapphire", "sapphire", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders(UserID, OrderID, Amount, CommitedAt) VALUES ("gold","ORDER20230912-072500Z", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO OrderDetails(UserID, OrderID, OrderDetailID, ItemID, Price, Quantity, CommitedAt) VALUES("gold", "ORDER20230912-072500Z", 1, "pen", 100, 1, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders(UserID, OrderID, Amount, CommitedAt) VALUES ("gold","ORDER20230912-082500Z", 400, PENDING_COMMIT_TIMESTAMP());
INSERT INTO OrderDetails(UserID, OrderID, OrderDetailID, ItemID, Price, Quantity, CommitedAt) VALUES("gold", "ORDER20230912-082500Z", 1, "pen", 100, 1, PENDING_COMMIT_TIMESTAMP());
INSERT INTO OrderDetails(UserID, OrderID, OrderDetailID, ItemID, Price, Quantity, CommitedAt) VALUES("gold", "ORDER20230912-082500Z", 2, "note", 150, 2, PENDING_COMMIT_TIMESTAMP());
```

以下のような感じで、指定したユーザのOrderを新しい順で取得するのが得意

```
EXPLAIN ANALYZE
SELECT
  *
FROM
  Orders
WHERE
  UserID = "ruby"
ORDER BY
  OrderID DESC
LIMIT 10
```

```
+----+---------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                      | Rows_Returned | Executions | Total_Latency |
+----+---------------------------------------------------------------------------+---------------+------------+---------------+
| *0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true) | 10            | 1          | 0.11 msecs    |
|  1 | +- Serialize Result                                                       | 10            | 1          | 0.08 msecs    |
|  2 |    +- Global Limit                                                        | 10            | 1          | 0.07 msecs    |
|  3 |       +- Local Distributed Union                                          | 10            | 1          | 0.07 msecs    |
| *4 |          +- Filter Scan                                                   |               |            |               |
|  5 |             +- Table Scan (Table: Orders, scan_method: Scalar)            | 10            | 1          | 0.06 msecs    |
+----+---------------------------------------------------------------------------+---------------+------------+---------------+
Predicates(identified by ID):
 0: Split Range: ($UserID = 'ruby')
 4: Seek Condition: ($UserID = 'ruby')

10 rows in set (4.54 msecs)
timestamp:            2023-09-12T20:47:11.536893+09:00
cpu time:             4.51 msecs
rows scanned:         10 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230911_11_44_36UTC
```
