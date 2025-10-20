CREATE TABLE Singers
(
    SingerId  INT64 NOT NULL,
    FirstName STRING(1024) NOT NULL,
    LastName  STRING(1024) NOT NULL,
) PRIMARY KEY (SingerId);

CREATE TABLE Albums
(
    SingerId INT64 NOT NULL,
    AlbumId  INT64 NOT NULL,
    Title    STRING(1024) NOT NULL,
) PRIMARY KEY (SingerId, AlbumId);

CREATE TABLE Users
(
    UserID    STRING(64) NOT NULL,
    UserName  STRING(1024) NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp= true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp= true),
) PRIMARY KEY (UserID);

CREATE TABLE UserBalances
(
    UserID    STRING(64) NOT NULL,
    Amount    INT64     NOT NULL,
    Point     INT64     NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp= true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp= true),
) PRIMARY KEY (UserID);

CREATE TABLE UserDepositHistories
(
    DepositID  STRING(64) NOT NULL,
    UserID     STRING(64) NOT NULL,
    Amount     INT64     NOT NULL,
    Point      INT64     NOT NULL,
    CommitedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp= true),
) PRIMARY KEY (DepositID);

CREATE INDEX UserIDByUserDepositHistories
    ON UserDepositHistories (
        UserID
    );

CREATE INDEX UserIDStoredAmountAndPointByUserDepositHistories
    ON UserDepositHistories (
        UserID
    ) STORING (
	    Amount,
	    Point
    );

CREATE TABLE Orders
(
    UserID     STRING(64) NOT NULL,
    OrderID    STRING(64) NOT NULL,
    Amount     INT64     NOT NULL,
    CommitedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp= true),
) PRIMARY KEY (UserID, OrderID),
  ROW DELETION POLICY (OLDER_THAN(CommitedAt, INTERVAL 90 DAY));

CREATE INDEX OrdersByUserIDAndCommitedAtDesc
    ON Orders (
        UserID,
        CommitedAt DESC
    );

CREATE TABLE OrderDetails
(
    OrderID       STRING(64) NOT NULL,
    OrderDetailID STRING(64) NOT NULL,
    ItemID        STRING(64) NOT NULL,
    Price         INT64     NOT NULL,
    Quantity      INT64     NOT NULL,
    CommitedAt    TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp= true),
) PRIMARY KEY (OrderID, OrderDetailID),
  ROW DELETION POLICY (OLDER_THAN(CommitedAt, INTERVAL 90 DAY));

CREATE INDEX ItemIDByOrderDetails
    ON OrderDetails (
        ItemID
    );
