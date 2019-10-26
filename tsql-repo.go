package echor

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Query struct
type Query struct {
	Where  string
	Offset int64
	Limit  int64
	Order  string
}

// TsqlRepoMetaInterface interface
type TsqlRepoMetaInterface interface {
	Name() string
	DatasourceID() string
}

// TsqlRepoCRUDInterface interface
type TsqlRepoCRUDInterface interface {
	Find() string
}

// TsqlRepoQuery func
func TsqlRepoQuery(tsqlRepoMetaInterface TsqlRepoMetaInterface, sqlStatement string, args ...interface{}) (*sql.Rows, error) {
	driver := getDatasourceDriver(tsqlRepoMetaInterface.DatasourceID(), "mysql")
	if driver != nil {
		db := driver.(*sql.DB)
		stmt, err := db.Prepare(sqlStatement)
		defer stmt.Close()
		if err != nil {
			fmt.Println("[TsqlRepoQuery error] ", err)
			fmt.Println("[TsqlRepoQuery error] ", sqlStatement)
			return nil, err
		}
		res, err := stmt.Query(args...)
		if err != nil {
			fmt.Println("[TsqlRepoQuery error] ", err)
			fmt.Println("[TsqlRepoQuery error] ", sqlStatement)
			return nil, err
		}
		return res, err
	}
	return nil, errors.New("INVALID_TSQL_DRIVER")
}

// TsqlRepoExec func
func TsqlRepoExec(tsqlRepoMetaInterface TsqlRepoMetaInterface, sqlStatement string, args ...interface{}) (int64, error) {
	driver := getDatasourceDriver(tsqlRepoMetaInterface.DatasourceID(), "mysql")
	if driver != nil {
		db := driver.(*sql.DB)
		stmt, err := db.Prepare(sqlStatement)
		defer stmt.Close()
		if err != nil {
			fmt.Println("[TsqlRepoQuery error] ", err)
			fmt.Println("[TsqlRepoQuery error] ", sqlStatement)
			return 0, err
		}
		res, err := stmt.Exec(args...)
		if err != nil {
			fmt.Println("[TsqlRepoQuery error] ", err)
			fmt.Println("[TsqlRepoQuery error] ", sqlStatement)
			return 0, err
		}
		return res.RowsAffected()
	}
	return 0, errors.New("INVALID_TSQL_DRIVER")
}

// TsqlRepoFind func
func TsqlRepoFind(tsqlRepoMetaInterface TsqlRepoMetaInterface, columns []string, query *Query, args ...interface{}) (rows *sql.Rows, err error) {
	name := tsqlRepoMetaInterface.Name()

	var limit int64 = 1000
	if query.Limit > 0 {
		limit = query.Limit
	}

	var offset int64
	if query.Offset > 0 {
		offset = query.Offset
	}

	c := "*"
	if len(columns) > 0 {
		c = strings.Join(columns, ",")
	}

	q := fmt.Sprintf("SELECT %s FROM `%s`", c, name)

	if len(query.Where) > 0 {
		q = fmt.Sprintf("%s WHERE %s", q, query.Where)
	}

	if query.Order != "" {
		q = fmt.Sprintf("%s ORDER BY %s", q, query.Order)
	}
	q = fmt.Sprintf("%s LIMIT %d OFFSET %d", q, limit, offset)

	if query.Where != "" {
		rows, err = TsqlRepoQuery(tsqlRepoMetaInterface, q, args...)
	} else {
		rows, err = TsqlRepoQuery(tsqlRepoMetaInterface, q)
	}
	return
}

// TsqlUpdate func
func TsqlUpdate(tsqlRepoMetaInterface TsqlRepoMetaInterface, cols []string, model interface{}, query *Query, args ...interface{}) (int64, error) {
	name := tsqlRepoMetaInterface.Name()

	if query == nil || query.Where == "" {
		return 0, errors.New("Unsafe update operation. Where should not empty")
	}
	setArr := make([]string, len(cols))
	for i, v := range cols {
		setArr[i] = fmt.Sprintf("%s=?", v)
	}
	setString := strings.Join(setArr, ",")
	var p []interface{}
	valueArr, err := tsqlStructProjectedArrValue(model, cols)
	if err != nil {
		return -1, err
	}
	p = append(valueArr, args...)
	q := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", name, setString, query.Where)
	return TsqlRepoExec(tsqlRepoMetaInterface, q, p...)
}

type insertModelOpts struct {
	skipIfConflict bool
}

// TsqlInsert func
func TsqlInsert(tsqlRepoMetaInterface TsqlRepoMetaInterface, model interface{}, opts *insertModelOpts) (int64, error) {
	name := tsqlRepoMetaInterface.Name()

	columns := tsqlStructFields(model)
	p := make([]string, len(columns))

	for i := range columns {
		p[i] = "?"
	}
	if opts == nil {
		opts = &insertModelOpts{
			skipIfConflict: false,
		}
	}
	var q string
	if opts.skipIfConflict {
		q = fmt.Sprintf("INSERT IGNORE INTO `%s` (%s) VALUES (%s)", name, strings.Join(columns, ","), strings.Join(p, ","))
	} else {
		q = fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", name, strings.Join(columns, ","), strings.Join(p, ","))
	}
	tmp, err := sqlStructArrValue(model)
	if err != nil {
		return -1, err
	}
	return TsqlRepoExec(tsqlRepoMetaInterface, q, tmp...)
}

// TsqlStructFields func
func tsqlStructFields(d interface{}) []string {
	t := reflect.TypeOf(d)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	columns := []string{}
	for i := 0; i < t.NumField(); i++ {
		columns = append(columns, t.Field(i).Tag.Get("json"))
	}
	return columns
}

// TsqlStructScan func
func tsqlStructScan(rows *sql.Rows, model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return errors.New("model must be pointer")
	}

	v = reflect.Indirect(v)
	t := v.Type()

	cols, _ := rows.Columns()

	var m map[string]interface{}

	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	if err := rows.Scan(columnPointers...); err != nil {
		return err
	}

	m = make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		m[colName] = *val
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i).Tag.Get("json")

		if item, ok := m[field]; ok {
			if v.Field(i).CanSet() {
				if item != nil {
					switch v.Field(i).Kind() {
					case reflect.String:
						s, ok := item.(string)
						if !ok {
							v.Field(i).SetString(string(item.([]uint8)))
						} else {
							v.Field(i).SetString(s)
						}
					case reflect.Float32, reflect.Float64:
						v.Field(i).SetFloat(item.(float64))
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						v.Field(i).SetInt(item.(int64))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						v.Field(i).SetUint(item.(uint64))
					default:
						return fmt.Errorf("%s field type not support yet %s", t.Field(i).Name, v.Field(i).Kind().String())
					}
				}
			}
		}
	}
	return nil
}

func sqlStructArrValue(d interface{}) ([]interface{}, error) {
	return tsqlStructProjectedArrValue(d, nil)
}

func tsqlStructProjectedArrValue(d interface{}, projectedColumns []string) ([]interface{}, error) {
	var tmp []interface{}
	v := reflect.ValueOf(d)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	projected := false
	flag := make(map[string]int)
	if projectedColumns != nil && len(projectedColumns) > 0 {
		projected = true
		for i, c := range projectedColumns {
			flag[c] = i
		}
		tmp = make([]interface{}, len(projectedColumns))
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("json")
		if projected {
			if val, ok := flag[tag]; ok {
				switch field.Kind() {
				case reflect.String:
					tmp[val] = field.Interface().(string)
				case reflect.Float32, reflect.Float64:
					tmp[val] = field.Interface().(float64)
				case reflect.Int:
					tmp[val] = field.Interface().(int)
				case reflect.Int8:
					tmp[val] = field.Interface().(int8)
				case reflect.Int16:
					tmp[val] = field.Interface().(int16)
				case reflect.Int64:
					tmp[val] = field.Interface().(int64)
				default:
					return nil, fmt.Errorf("structValue (projected) - Type not support")
				}
			}
		} else {
			switch field.Kind() {
			case reflect.String:
				tmp = append(tmp, field.Interface().(string))
			case reflect.Float32, reflect.Float64:
				tmp = append(tmp, field.Interface().(float64))
			case reflect.Int:
				tmp = append(tmp, field.Interface().(int))
			case reflect.Int8:
				tmp = append(tmp, field.Interface().(int8))
			case reflect.Int16:
				tmp = append(tmp, field.Interface().(int16))
			case reflect.Int32:
				tmp = append(tmp, field.Interface().(int32))
			case reflect.Int64:
				tmp = append(tmp, field.Interface().(int64))
			default:
				return nil, fmt.Errorf("structValue - Type not support")
			}
		}
	}
	return tmp, nil
}
