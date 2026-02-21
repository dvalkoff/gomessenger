package repository

import "database/sql"

type WorkspaceRow struct {
	Id        int
	URL       string
	IsCurrent bool
}

type WorkspaceRepository interface {
	GetWorkspaces() ([]WorkspaceRow, error)
	GetCurrentSpace() (*WorkspaceRow, error)
	AddWorkspace(WorkspaceRow) (WorkspaceRow, error)
	SetWorkspaceCurrent(id int) error
}

type workspaceRepository struct {
	db *sql.DB
}

func NewWorkspaceRepository(db *sql.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) GetWorkspaces() ([]WorkspaceRow, error) {
	sql := `
	SELECT id, url, is_current FROM workspace
	ORDER BY id
	`
	rows, err := r.db.Query(sql)
	if err != nil {
		return nil, err
	}
	workspaceRows := make([]WorkspaceRow, 0)
	for rows.Next() {
		row := WorkspaceRow{}
		rows.Scan(&row.Id, &row.URL, &row.IsCurrent)
		workspaceRows = append(workspaceRows, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workspaceRows, nil
}

func (r *workspaceRepository) GetCurrentSpace() (*WorkspaceRow, error) {
	workspaces, err := r.GetWorkspaces()
	if err != nil {
		return nil, err
	}
	var currentWorkspace *WorkspaceRow
	for _, workspace := range workspaces {
		if workspace.IsCurrent {
			currentWorkspace = &workspace
			break
		}
	}
	return currentWorkspace, nil
}

func (r *workspaceRepository) AddWorkspace(workspace WorkspaceRow) (WorkspaceRow, error) {
	sql := `
	INSERT INTO workspace(url) VALUES($1);
	`
	result, err := r.db.Exec(sql, workspace.URL)
	if err != nil {
		return WorkspaceRow{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return WorkspaceRow{}, err
	}
	workspace.Id = int(id)
	return workspace, nil
}

func (r *workspaceRepository) SetWorkspaceCurrent(id int) error {
	sqlResetCurrent := `
	UPDATE workspace
	SET is_current = false
	`
	if _, err := r.db.Exec(sqlResetCurrent); err != nil {
		return err
	}
	sqlSetCurrent := `
	UPDATE workspace
	SET is_current = true
	WHERE id = $1
	`
	_, err := r.db.Exec(sqlSetCurrent, id)
	return err
}
