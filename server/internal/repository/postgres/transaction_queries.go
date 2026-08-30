package postgres

const CreateTransactionQuery = `
INSERT INTO transactions (id, user_id, category_id, amount, type, note, source, transaction_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING created_at, updated_at
`
