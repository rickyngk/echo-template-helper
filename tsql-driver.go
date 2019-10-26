package echor

import (
	"database/sql"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// TsqlDriver class
type TsqlDriver struct {
	connector *sql.DB
}

// NewTsqlDriver - create new connection
func NewTsqlDriver(conf DatasourceSpecs) DbDriverInterface {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Username, conf.Password, conf.Address, conf.Port, conf.DbName)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	driver := &TsqlDriver{connector: db}
	return driver
}

// GetNativeConnector func - impl interface DbNativeInterface
func (d *TsqlDriver) GetNativeConnector() interface{} {
	return d.connector
}
