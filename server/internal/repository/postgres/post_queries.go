package postgres

const queryCreatePost = `
	INSERT INTO posts (id, user_id, image_url, caption, visibility, location_name, transaction_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
`