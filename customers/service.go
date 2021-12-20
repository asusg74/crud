package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	log.Print("EnteR")

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

func (s *Service) AllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE active = TRUE
	`)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return items, nil
}
func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, phone, active, created FROM customers
	`)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return items, nil
}
func (s *Service) Block(ctx context.Context, id int64) error {

	_, err := s.ByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return ErrInternal
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE customers SET active = FALSE WHERE id = $1
	`, id)
	if err != nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}
func (s *Service) Unblock(ctx context.Context, id int64) error {

	_, err := s.ByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return ErrInternal
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE customers SET active = TRUE WHERE id = $1
	`, id)
	if err != nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}
func (s *Service) Remove(ctx context.Context, id int64) error {

	_, err := s.ByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return ErrInternal
	}
	_, err = s.db.ExecContext(ctx, `
		DELETE FROM customers WHERE id = $1
	`, id)
	if err != nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}
func (s *Service) Save(ctx context.Context, id int64, phone string, name string) (*Customer, error) {
	result := &Customer{}
	if id == 0 {
		err := s.db.QueryRowContext(ctx, `
			INSERT INTO customers (name, phone) VALUES($1, $2) RETURNING id, name, phone, active, created
		`, name, phone).Scan(&result.ID, &result.Name, &result.Phone, &result.Active, &result.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
	} else {
		_, err := s.ByID(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		if err != nil {
			return nil, ErrInternal
		}
		_, err = s.db.ExecContext(ctx, `
			UPDATE customers SET phone = $2, name = $3 WHERE id = $1
		`, id, phone, name)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		result, err = s.ByID(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		if err != nil {
			return nil, ErrInternal
		}
	}

	return result, nil
}
