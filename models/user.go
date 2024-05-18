package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

// FindUserLists finds all the todo lists associated to the user. It also
// finds all associated todo for every list
func FindUserLists(dao *daos.Dao, id string) ([]*List, error) {
	lists := []*List{}

	if err := ListQuery(dao).AndWhere(dbx.HashExp{"owner": id}).All(&lists); err != nil {
		return nil, err
	}

	for _, l := range lists {
		if err := TaskQuery(dao).AndWhere(dbx.HashExp{"list": l.Id}).All(&l.Tasks); err != nil {
			return nil, err
		}

		l.Completed = 0
		for _, t := range l.Tasks {
			if t.Done {
				l.Completed++
			}
		}
	}

	return lists, nil
}
