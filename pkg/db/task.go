package db

import (
	"database/sql"
	"time"
)

const TaskLimit = 50

type Task struct {
	ID      int64  `db:"id" json:"id,string"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
INSERT INTO scheduler (date, title, comment, repeat)
VALUES (?, ?, ?, ?)
`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = TaskLimit
	}

	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat
         FROM scheduler
         ORDER BY date
         LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func UpdateTaskDate(id string, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	result, err := DB.Exec(query, date, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// SearchTasks ищет задачи по тексту или дате
func SearchTasks(search string) ([]*Task, error) {
	var tasks []*Task

	// Пытаемся распарсить как дату в формате DD.MM.YYYY
	dateFormat := "02.01.2006"
	parsedDate, err := time.Parse(dateFormat, search)

	var rows *sql.Rows

	if err == nil {
		// Это дата — ищем по полю date
		dateStr := parsedDate.Format("20060102")
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat 
             FROM scheduler 
             WHERE date = ? 
             ORDER BY date 
             LIMIT ?`, dateStr, TaskLimit)
		if err != nil {
			return nil, err
		}
	} else {
		// Это текст — ищем по title и comment
		searchPattern := "%" + search + "%"
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat 
             FROM scheduler 
             WHERE title LIKE ? OR comment LIKE ? 
             ORDER BY date 
             LIMIT ?`, searchPattern, searchPattern, TaskLimit)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}
