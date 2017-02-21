package store

import "database/sql"

// Task ...
type Task struct {
	ID          int
	Description string
	CreatorID   int
}

const sqlInsertTask = `INSERT INTO tasks (description, creator_id) ` +
	`VALUES ($1, $2) RETURNING id`

// CreateTask ...
func CreateTask(db *sql.DB, descr string, creatorID int) (*Task, error) {
	t := &Task{Description: descr, CreatorID: creatorID}
	row := db.QueryRow(sqlInsertTask, t.Description, t.CreatorID)
	err := row.Scan(&t.ID)
	return t, err
}

const sqlSelectTasks = `SELECT id, description, creator_id FROM tasks ` +
	`WHERE creator_id = $1`

// ListTasks ...
func ListTasks(db *sql.DB, userID int) ([]Task, error) {
	rows, err := db.Query(sqlSelectTasks, userID)
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
