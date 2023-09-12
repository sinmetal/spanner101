CREATE TABLE Orders (
    UserID STRING(64) NOT NULL,
    OrderID STRING(64) NOT NULL,
    Amount INT64 NOT NULL,
    CommitedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (UserID, OrderID), INTERLEAVE IN PARENT Users ON DELETE CASCADE,
  ROW DELETION POLICY (OLDER_THAN(CommitedAt, INTERVAL 90 DAY));

CREATE INDEX UserIDAndCommitedAtDescByOrdersParentUsers
ON Orders (
    UserID,
    CommitedAt DESC
), INTERLEAVE IN Users;

CREATE TABLE OrderDetails (
    UserID STRING(64) NOT NULL,
    OrderID STRING(64) NOT NULL,
    OrderDetailID STRING(64) NOT NULL,
    ItemID STRING(64) NOT NULL,
    Price INT64 NOT NULL,
    Quantity INT64 NOT NULL,
    CommitedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (UserID, OrderID, OrderDetailID), INTERLEAVE IN PARENT Orders ON DELETE CASCADE;

CREATE INDEX ItemIDByOrderDetails
ON OrderDetails (
    ItemID
);
