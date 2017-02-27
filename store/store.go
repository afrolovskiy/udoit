package store

import (
	"database/sql"
	"fmt"
)

// Task ...
type Task struct {
	ID          int
	IDInChat    int
	ChatID      int64
	Description string
	CreatorID   int
	AssigneeID  int
}

const sqlCreateSequence = "CREATE SEQUENCE IF NOT EXISTS "
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
	db.Query(sqlCreateSequence + "\"" + seqname + "\"")

	t := &Task{Description: desc, CreatorID: creatorID, ChatID: chatID}

	rowSeq := db.QueryRow(sqlSelectNextSeqVal, seqname)
	rowSeq.Scan(&t.IDInChat)

	rowTask := db.QueryRow(sqlInsertTask, t.IDInChat, t.ChatID, t.CreatorID, t.Description)
	err := rowTask.Scan(&t.ID)
	return t, err
}

const sqlDeleteTask = "DELETE FROM tasks WHERE chat_id = $1 AND id_in_chat = $2"

// DeleteTask ...
func DeleteTask(db *sql.DB, chatID int64, ID int) error {
	_, err := db.Exec(sqlDeleteTask, chatID, ID)
	return err
}

const sqlSelectTasks = `SELECT id_in_chat, description, creator_id FROM tasks ` +
	`WHERE chat_id = $1`

// ChatTasks ...
func ChatTasks(db *sql.DB, chatID int64) ([]Task, error) {
	rows, err := db.Query(sqlSelectTasks, chatID)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for rows.Next() {
		t := Task{}
		err = rows.Scan(&t.IDInChat, &t.Description, &t.CreatorID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	err = rows.Err()
	return tasks, err
}

const sqlSelectUserTasks = sqlSelectTasks + " OR assignee_id = $2"

// UserTasks ...
func UserTasks(db *sql.DB, userID int) ([]Task, error) {
	rows, err := db.Query(sqlSelectUserTasks, userID, userID)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for rows.Next() {
		t := Task{}
		err = rows.Scan(&t.IDInChat, &t.Description, &t.CreatorID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
