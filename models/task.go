package models

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Task struct {
	models.BaseModel

	// Fields
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Done        bool           `db:"done"`
	Deadline    types.DateTime `db:"deadline"`

	// Relations
	ListId string `db:"list"`
	List   *List  `db:"-"`
}

func (t *Task) Save(dao *daos.Dao) error {
	if err := dao.Save(t); err != nil {
		return fmt.Errorf("save task (%s): %w", t.Id, err)
	}

	return nil
}

func (*Task) TableName() string {
	return "tasks"
}

func TaskQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Task{})
}

func GetTaskById(dao *daos.Dao, id string) (*Task, error) {
	t := &Task{}
	tq := dao.ModelQuery(t)

	if err := tq.AndWhere(dbx.HashExp{"id": id}).Limit(1).One(t); err != nil {
		return nil, fmt.Errorf("db query task (%s): %w", id, err)
	}

	return t, nil
}
