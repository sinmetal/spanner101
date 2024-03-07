CREATE INDEX UserIDAndCommitedAtDescStoringAmountParentUsersByOrders
ON Orders (
    UserID,
    CommitedAt DESC
) STORING (Amount), INTERLEAVE IN Users;