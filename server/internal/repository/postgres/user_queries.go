package postgres


const (
	queryCreateUser = `
		INSERT INTO users (id, username, email, password_hash, full_name, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	queryGetByEmail = `
		SELECT id, username, email, password_hash, status, created_at
		FROM users
		WHERE email = $1 
		LIMIT 1
	`
)