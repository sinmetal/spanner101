# GROUP BY pattern1

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
INSERT INTO Orders (OrderID, UserID, Amount, CommitedAt) VALUES ("00000005-0543-420d-ae8e-d00cb1c99cc1", "ruby", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (OrderID, UserID, Amount, CommitedAt) VALUES ("00000694-1c17-4812-93c0-06070df608f5", "ruby", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (OrderID, UserID, Amount, CommitedAt) VALUES ("0001b069-5bbf-4397-bc93-57d6252b1b17", "dia", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (OrderID, UserID, Amount, CommitedAt) VALUES ("00054374-14b6-4bdd-9a04-6b39187465fc", "sapphire", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (OrderID, UserID, Amount, CommitedAt) VALUES ("0008c6e2-8278-4d23-9270-4c03a4d54534", "silver", 100, PENDING_COMMIT_TIMESTAMP());
INSERT INTO Orders (OrderID, UserID, Amount, CommitedAt) VALUES ("000cf780-bdf5-4932-9aac-1fb6cfb893e2", "gold", 100, PENDING_COMMIT_TIMESTAMP());
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
+----+--------------------------------------------------------------------------------------+---------------+------------+---------------+
| ID | Query_Execution_Plan                                                                 | Rows_Returned | Executions | Total_Latency |
+----+--------------------------------------------------------------------------------------+---------------+------------+---------------+
|  0 | Serialize Result                                                                     | 5             | 1          | 6.2 msecs     |
|  1 | +- Global Stream Aggregate                                                           | 5             | 1          | 6.2 msecs     |
|  2 |    +- Sort                                                                           | 5             | 1          | 6.19 msecs    |
|  3 |       +- Distributed Union (distribution_table: Orders, split_ranges_aligned: false) | 5             | 1          | 6.18 msecs    |
|  4 |          +- Local Hash Aggregate                                                     | 5             | 1          | 6.15 msecs    |
|  5 |             +- Local Distributed Union                                               | 4419          | 1          | 5.63 msecs    |
|  6 |                +- Table Scan (Full scan: true, Table: Orders, scan_method: Scalar)   | 4419          | 1          | 5.37 msecs    |
+----+--------------------------------------------------------------------------------------+---------------+------------+---------------+
5 rows in set (10.55 msecs)
timestamp:            2023-09-15T14:38:43.309677+09:00
cpu time:             8.39 msecs
rows scanned:         4419 rows
deleted rows scanned: 0 rows
optimizer version:    5
optimizer statistics: auto_20230909_06_42_57UTC
```

GROUP BYのために `Local Hash Aggregate` が登場。
この例だとUserIDの種類が5つしかないので、すぐ終わっているが、Userが多いとHash Table作って集計するのが大変になる。
Localでの集計が終わった後、 `Sort` して、 `Global Stream Aggregate` で最終集計を行っている。