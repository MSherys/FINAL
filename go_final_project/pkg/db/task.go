package dbase

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// ************************************************
func AddTask(task *Task) (int64, error) {
	var id int64
	if DB == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ************************************************
func Tasks(limit int) ([]*Task, error) {
	var res []*Task
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// ************************************************
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :task_id",
		sql.Named("task_id", id))

	task := &Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// ****************************************************************
func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	row, err := DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :task_id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("task_id", task.ID))
	if err != nil {
		return err
	}
	count, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// ***************************************************
func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	result, err := DB.Exec("DELETE FROM scheduler WHERE id = :task_id",
		sql.Named("task_id", id))
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("запись не найдена")
	}
	return nil
}

// ***************************************************

func UpdateDate(next string, id string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	row, err := DB.Exec("UPDATE scheduler SET date = :new_date WHERE id = :task_id",
		sql.Named("new_date", next),
		sql.Named("task_id", id))
	if err != nil {
		return err
	}
	count, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}
