package postgres


const (
	queryCreateUser = `
		INSERT INTO users (id,  email, password_hash,username, firstname, lastname, phone_number, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	queryGetByEmail = `
		SELECT id, email, firstname, lastname,  password_hash,username, avatar_url, phone_number,bio, status, created_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`
)