package store

import (
	"database/sql"
	"fmt"
)

// Task ...
type Task struct {
	ID          int
	IDinchat    int
	ChatID      int64
	Description string
	CreatorID   int
}

const sqlCreateSequence = "CREATE SEQUENCE IF NOT EXISTS $1"
const sqlSelectNextSeqVal = "SELECT nextval($1)"
const sqlInsertTask = `INSERT INTO tasks (id_in_chat, chat_id, creator_id, description) ` +
	`VALUES ($1, $2, $3, $4) RETURNING id`

func createSequenceName(chatID int64) string {
	name := "chat_" + fmt.Sprintf("%d", chatID) + "_seq"
	return name
}

// CreateTask ...
func CreateTask(db *sql.DB, desc string, creatorID int, chatID int64) (*Task, error) {
	seqname := createSequenceName(chatID)
	db.Query(sqlCreateSequence, seqname)

	t := &Task{Description: desc, CreatorID: creatorID, ChatID: chatID}

	rowSeq := db.QueryRow(sqlSelectNextSeqVal, seqname)
	rowSeq.Scan(&t.IDinchat)

	rowTask := db.QueryRow(sqlInsertTask, t.IDinchat, t.ChatID, t.CreatorID, t.Description)
	err := rowTask.Scan(&t.ID)
	return t, err
}

const sqlSelectTasks = `SELECT id_in_chat, description, creator_id FROM tasks ` +
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
		err = rows.Scan(&t.IDinchat, &t.Description, &t.CreatorID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	err = rows.Err()
	return tasks, err
}
