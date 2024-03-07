CREATE INDEX UserIDAndCommitedAtDescStoringAmountByOrders
ON Orders (
    UserID,
    CommitedAt DESC
) STORING (Amount);