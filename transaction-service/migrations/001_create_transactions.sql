CREATE TABLE transactions (
                              id TEXT PRIMARY KEY,
                              purchase_id TEXT,
                              transaction_id TEXT,
                              amount BIGINT,
                              status TEXT
);