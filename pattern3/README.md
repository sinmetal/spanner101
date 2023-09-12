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
