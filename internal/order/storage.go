package order

import (
	"database/sql"
)

type OrderStorage struct {
	db *sql.DB
}

func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{
		db: db,
	}
}

func (s *OrderStorage) createOrder(title, description string, completed bool) (string, error) {
	order := Order{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
	statement := `insert into orders(title, description, completed) values($1, $2, $3);`

	_, err := s.db.Exec(statement, order.Title, order.Description, order.Completed)
	if err != nil {
		return "creation had an error", err
	}

	return "Successfully created", nil
}

func (s *OrderStorage) getAllOrders() ([]Order, error) {
	var orders []Order
	statement := `select * from orders;`
	rows, err := s.db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Title, &order.Description, &order.Completed)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
