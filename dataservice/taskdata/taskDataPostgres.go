package taskdata

import (
	"database/sql"
	"scarlet_backend/model"
)

type TaskDataPostgres struct{ db *sql.DB }

func NewTaskDataPostgres(db *sql.DB) *TaskDataPostgres { return &TaskDataPostgres{db: db} }

func (t *TaskDataPostgres) GetTask() ([]model.Task, error) {
	rows, err := t.db.Query("SELECT taskId, team_id, fire_id, description, status, active, created_at, created_by FROM task_data.tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err = rows.Scan(&task.ID, &task.TeamID, &task.FireID, &task.Description, &task.Status, &task.Active, &task.CreatedAt, &task.CreatedBy); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *TaskDataPostgres) AddTask(task model.Task) error {
	query := `insert into task_data.tasks(team_id, fire_id, description, status, active, created_at, created_by) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := t.db.Exec(query, task.TeamID, task.FireID, task.Description, task.Status, task.Active, task.CreatedAt, task.CreatedBy)
	return err
}
