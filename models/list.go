package models

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type List struct {
	models.BaseModel

	// Fields
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Archived    bool           `db:"archived"`
	Deadline    types.DateTime `db:"deadline"`

	// Relations
	OwnerId string         `db:"owner"`
	Owner   *models.Record `db:"-"`
	Tasks   []*Task        `db:"-"`

	// Computed
	Completed int `db:"-"`
}

func (l *List) Save(dao *daos.Dao) error {
	if err := dao.Save(l); err != nil {
		return fmt.Errorf("save list (%s): %w", l.Id, err)
	}
	return nil
}

func (*List) TableName() string {
	return "lists"
}

func ListQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&List{})
}

func GetListById(dao *daos.Dao, id string, expand bool) (*List, error) {
	l := &List{}

	if err := ListQuery(dao).AndWhere(dbx.HashExp{"id": id}).Limit(1).One(&l); err != nil {
		return nil, err
	}

	if err := TaskQuery(dao).AndWhere(dbx.HashExp{"list": l.Id}).All(&l.Tasks); err != nil {
		return nil, err
	}

	l.Completed = 0
	for _, t := range l.Tasks {
		if t.Done {
			l.Completed++
		}
	}

	return l, nil
}
