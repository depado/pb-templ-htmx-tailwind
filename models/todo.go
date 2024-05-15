package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ToDoList struct {
	models.BaseModel

	// Fields
	Title string `db:"title"`

	// Relations
	Owner string `db:"owner"`
}

func (*ToDoList) TableName() string {
	return "todo_lists"
}

type ToDo struct {
	models.BaseModel

	// Fields
	Name        string         `db:"title"`
	Description string         `db:"description"`
	Done        bool           `db:"done"`
	Deadline    types.DateTime `db:"deadline"`

	// Relations
	List  string `db:"list"`
	Owner string `db:"owner"`
}

func (*ToDo) TableName() string {
	return "todos"
}
