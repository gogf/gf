package test

import (
	"database/sql"

	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/frame/gins"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
)

type User struct {
	Id         int         `orm:"id,primary"  json:"id"`
	Passport   string      `orm:"passport"    json:"passport"`
	Password   string      `orm:"password"    json:"password"`
	NickName   string      `orm:"nickname"    json:"nick_name"`
	CreateTime *gtime.Time `orm:"create_time" json:"create_time"`
}

type UserModel struct {
	*gdb.Model
	TableName string
}

var (
	UserTableName      = "user"
	gUserModelCacheKey = gdebug.CallerFilePath()
)

func ModelUser() *UserModel {
	return gins.GetOrSetFunc(gUserModelCacheKey, func() interface{} {
		return &UserModel{
			DB().Table(UserTableName).Safe(),
			UserTableName,
		}
	}).(*UserModel)
}

func (r *User) Insert() (result sql.Result, err error) {
	return ModelUser().Data(r).Insert()
}

func (r *User) Replace() (result sql.Result, err error) {
	return ModelUser().Data(r).Replace()
}

func (r *User) Save() (result sql.Result, err error) {
	return ModelUser().Data(r).Save()
}

func (r *User) Update() (result sql.Result, err error) {
	return ModelUser().Data(r).Where(gdb.GetWhereConditionOfStruct(r)).Update()
}

func (r *User) Delete() (result sql.Result, err error) {
	return ModelUser().Where(gdb.GetWhereConditionOfStruct(r)).Delete()
}

func (m *UserModel) Select() ([]*User, error) {
	return m.All()
}

func (m *UserModel) All() ([]*User, error) {
	array := ([]*User)(nil)
	if err := m.Scan(&array); err != nil {
		return nil, err
	}
	return array, nil
}

func (m *UserModel) One() (*User, error) {
	list, err := m.All()
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}
