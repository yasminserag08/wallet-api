# How the Transfer Endpoint Prevents Race Conditions

## The Problem

Consider two simultaneous withdrawal requests against a wallet with a balance of 100:

1. Request A reads balance: 100
2. Request B reads balance: 100
3. Request A checks: 100 >= 100, proceeds
4. Request B checks: 100 >= 100, proceeds
5. Request A deducts 100, balance becomes 0
6. Request B deducts 100, balance becomes -100

Both requests passed the balance check because they both read the balance before either
one updated it. This is a classic race condition.

## The Solution: Database Transactions + Row-Level Locking

The transfer endpoint (and all balance-changing operations) use two mechanisms together:

### 1. Database Transactions
All steps of a transfer run inside a single GORM database transaction:
- Deduct from sender
- Credit receiver
- Create transfer_out record
- Create transfer_in record

If any step fails, the entire transaction rolls back, which guarantees atomicity.

### 2. Row-Level Locking (SELECT FOR UPDATE)
Before reading a wallet's balance, the row is locked:

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).
    Where("user_id = ?", userID).
    First(&wallet)
```

This translates to `SELECT * FROM wallets WHERE user_id = ? FOR UPDATE` in SQL.

When Request A locks the wallet row, Request B is blocked from reading that same row
until Request A's transaction commits or rolls back. This forces the requests to queue:

1. Request A locks wallet row, reads balance: 100
2. Request B tries to lock the same row — blocked, must wait
3. Request A deducts 100, balance becomes 0, commits, releases lock
4. Request B gets the lock, reads balance: 0
5. Request B checks: 0 >= 100, returns insufficient funds

Only one withdrawal succeeds. The balance never goes negative.

## Why Both Are Needed

- Transactions alone don't prevent the race condition — two transactions can still
  read the same balance simultaneously before either commits.
- Row locking alone without a transaction would release the lock too early (lock is immediately released after `SELECT` before the `UPDATE` even runs).