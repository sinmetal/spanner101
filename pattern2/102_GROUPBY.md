# GROUP BY pattern2

UserごとにAmountを集計するサンプル。

# Sample Data

```
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("ruby", "ruby", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("dia", "dia", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("sapphire", "sapphire", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("silver", "silver", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
INSERT INTO Users (UserID, UserName, CreatedAt, UpdatedAt) VALUES ("gold", "gold", PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
```

```
INSERT INTO Orders (UserID, OrderID, Amount, CommitedAt) VALUES ("ruby", "00000005-0543-420d-ae8e-d00cb1c99cc1", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (UserID, OrderID, Amount, CommitedAt) VALUES ("ruby", "00000694-1c17-4812-93c0-06070df608f5", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (UserID, OrderID, Amount, CommitedAt) VALUES ("dia", "0001b069-5bbf-4397-bc93-57d6252b1b17", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (UserID, OrderID, Amount, CommitedAt) VALUES ("sapphire", "00054374-14b6-4bdd-9a04-6b39187465fc", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (UserID, OrderID, Amount, CommitedAt) VALUES ("silver", "0008c6e2-8278-4d23-9270-4c03a4d54534", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (UserID, OrderID, Amount, CommitedAt) VALUES ("gold", "000cf780-bdf5-4932-9aac-1fb6cfb893e2", 100, PENDING_COMMIT_TIMESTAMP());
```

```
EXPLAIN ANALYZE
SELECT
  UserID,
  SUM(Orders.Amount) AS Amount
FROM
  Orders
GROUP BY
  UserID
```

```
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                         | Rows_Returned | Executions | Total_Latency |
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Distributed Union (distribution_table: Users, split_ranges_aligned: true)    | 5             | 1          | 3.15 msecs    |
|  1 | +- Serialize Result                                                          | 5             | 1          | 3.13 msecs    |
|  2 |    +- Stream Aggregate                                                       | 5             | 1          | 3.13 msecs    |
|  3 |       +- Local Distributed Union                                             | 4353          | 1          | 2.86 msecs    |
|  4 |          +- Table Scan (Full scan: true, Table: Orders, scan_method: Scalar) | 4353          | 1          | 2.6 msecs     |
+----+------------------------------------------------------------------------------+---------------+------------+---------------+
5 rows in set (7.53 msecs)
timestamp:            2023-09-15T14:35:53.678128+09:00
cpu time:             6.77 msecs
rows scanned:         4353 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230906_07_18_51UTC
```

PKの先頭にUserIDがあるので、UserごとにOrderはまとまっている。
そのため、Stream Aggregateで集計ができる。
Orders TableはUsers Tableの子どもなので、Localで完結している。