package test

import (
	"database/sql"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id         int         `orm:"id,primary"  json:"id"`
	Passport   string      `orm:"passport"    json:"passport"`
	Password   string      `orm:"password"    json:"password"`
	NickName   string      `orm:"nickname"    json:"nick_name"`
	CreateTime *gtime.Time `orm:"create_time" json:"create_time"`
}

// UserModel is the model of convenient operations for table user.
type UserModel struct {
	*gdb.Model
	TableName string
}

var (
	// UserTableName is the table name of user.
	UserTableName = "user"
)

// ModelUser creates and returns a new model object for table user.
func ModelUser() *UserModel {
	return &UserModel{
		g.DB(ConfigGroup).Table(UserTableName).Safe(),
		UserTableName,
	}
}

// Inserts does "INSERT...INTO..." statement for inserting current object into table.
func (r *User) Insert() (result sql.Result, err error) {
	return ModelUser().Data(r).Insert()
}

// Replace does "REPLACE...INTO..." statement for inserting current object into table.
// If there's already another same record in the table (it checks using primary key or unique index),
// it deletes it and insert this one.
func (r *User) Replace() (result sql.Result, err error) {
	return ModelUser().Data(r).Replace()
}

// Save does "INSERT...INTO..." statement for inserting/updating current object into table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
func (r *User) Save() (result sql.Result, err error) {
	return ModelUser().Data(r).Save()
}

// Update does "UPDATE...WHERE..." statement for updating current object from table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
func (r *User) Update() (result sql.Result, err error) {
	return ModelUser().Data(r).Where(gdb.GetWhereConditionOfStruct(r)).Update()
}

// Delete does "DELETE FROM...WHERE..." statement for deleting current object from table.
func (r *User) Delete() (result sql.Result, err error) {
	return ModelUser().Where(gdb.GetWhereConditionOfStruct(r)).Delete()
}

// Select overwrite the Select method from gdb.Model for model
// as retuning all objects with specified structure.
func (m *UserModel) Select() ([]*User, error) {
	array := ([]*User)(nil)
	if err := m.Scan(&array); err != nil {
		return nil, err
	}
	return array, nil
}

// First does the same logistics as One method from gdb.Model for model
// as retuning first/one object with specified structure.
func (m *UserModel) First() (*User, error) {
	list, err := m.Select()
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}
