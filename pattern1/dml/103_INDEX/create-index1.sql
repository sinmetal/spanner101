CREATE INDEX OrdersByUserIDAndCommitedAtDescStoringAmount
ON Orders (
    UserID,
    CommitedAt DESC
) STORING (Amount);