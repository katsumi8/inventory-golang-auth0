package todo

import (
	"database/sql"
)

type TodoStorage struct {
	db *sql.DB
}

func NewTodoStorage(db *sql.DB) *TodoStorage {
	return &TodoStorage{
		db: db,
	}
}

func (s *TodoStorage) createTodo(title, description string, completed bool) (string, error) {
	todo := Todo{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
	statement := `insert into todos(title, description, completed) values($1, $2, $3);`

	_, err := s.db.Exec(statement, todo.Title, todo.Description, todo.Completed)
	if err != nil {
		return "creation had an error", err
	}

	return "Successfully created", nil
}

func (s *TodoStorage) getAllTodos() ([]Todo, error) {
	var todos []Todo
	statement := `select * from todos;`
	rows, err := s.db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt,
			&todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
