package repository

import (
	"context"
	"database/sql"
	"errors"
	"template/internal/phonebook"
)

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db}
}

func (r *PostgreSQLRepository) NewUser(ctx context.Context, user *phonebook.User) (int, error) {
	var id int

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING ID`,
		user.Email,
		user.Password,
	).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *PostgreSQLRepository) GetUserByEmail(ctx context.Context, email string) (*phonebook.User, error) {
	var user phonebook.User

	err := r.db.QueryRowContext(ctx, `SELECT id, email, password FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgreSQLRepository) NewAddress(ctx context.Context, address *phonebook.Address) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO addresses (user_id, name, phone_number) VALUES ($1, $2, $3)`,
		address.User.ID,
		address.Name,
		address.PhoneNumber,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgreSQLRepository) Addresses(ctx context.Context) ([]*phonebook.Address, error) {
	res := make([]*phonebook.Address, 0)

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, phone_number FROM addresses`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cur := phonebook.Address{User: &phonebook.User{}}
		err = rows.Scan(
			&cur.ID,
			&cur.User.ID,
			&cur.Name,
			&cur.PhoneNumber,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, &cur)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *PostgreSQLRepository) GetAddressesByUserID(ctx context.Context, userID int) ([]*phonebook.Address, error) {
	res := make([]*phonebook.Address, 0)

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, phone_number FROM addresses WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cur := phonebook.Address{User: &phonebook.User{}}
		err = rows.Scan(
			&cur.ID,
			&cur.User.ID,
			&cur.Name,
			&cur.PhoneNumber,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, &cur)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *PostgreSQLRepository) GetAddressByID(ctx context.Context, ID int) (*phonebook.Address, error) {
	address := phonebook.Address{User: &phonebook.User{}}

	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, phone_number FROM addresses WHERE id = $1`, ID).
		Scan(&address.ID, &address.User.ID, &address.Name, &address.PhoneNumber)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *PostgreSQLRepository) UpdateAddress(ctx context.Context, ID int, address *phonebook.Address) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE addresses SET name = $1, phone_number = $2 WHERE id = $3`,
		address.Name,
		address.PhoneNumber,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgreSQLRepository) DeleteAddress(ctx context.Context, ID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM addresses WHERE id = $1`, ID)
	if err != nil {
		return err
	}

	return nil
}
