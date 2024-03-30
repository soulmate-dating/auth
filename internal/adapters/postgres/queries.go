package postgres

const (
	getUserByEmailQuery = `SELECT * FROM auth.users WHERE email = $1`
	createUserQuery     = `INSERT INTO auth.users (
                      		id, email, password
    						) VALUES ($1, $2, $3) RETURNING id`
	//updateUserLoginStatusQuery = `UPDATE users SET logged_in = $2 WHERE id = $1 RETURNING *`
	//updatePasswordQuery = `UPDATE users password = $2 WHERE id = $1 RETURNING *`
)
