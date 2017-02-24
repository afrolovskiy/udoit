package store

import "database/sql"

// Task ...
type Task struct {
	ID          int
	ChatID      int64
	Description string
	CreatorID   int
}

const sqlInsertTask = `INSERT INTO tasks (description, creator_id, chat_id) ` +
	`VALUES ($1, $2, $3) RETURNING id`

// CreateTask ...
func CreateTask(db *sql.DB, desc string, creatorID int, chatID int64) (*Task, error) {
	t := &Task{Description: desc, CreatorID: creatorID, ChatID: chatID}
	row := db.QueryRow(sqlInsertTask, t.Description, t.CreatorID, t.ChatID)
	err := row.Scan(&t.ID)
	return t, err
}

const sqlSelectTasks = `SELECT id, description, creator_id FROM tasks ` +
	`WHERE chat_id = $1`

// ListTasks ...
func ListTasks(db *sql.DB, chatID int64) ([]Task, error) {
	rows, err := db.Query(sqlSelectTasks, chatID)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for rows.Next() {
		t := Task{}
		err = rows.Scan(&t.ID, &t.Description, &t.CreatorID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	err = rows.Err()
	return tasks, err
}
