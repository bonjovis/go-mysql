package mysql

//mysql
import "fmt"
import "time"
import "strconv"
import "strings"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type DbPool struct {
	db *sql.DB
}

func NewDatabaseConnectionPool(user, pwd, host string) *DbPool {
	db, _ := sql.Open("mysql", user+":"+pwd+"@"+host)
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	dbPool := DbPool{db}
	return &dbPool
}

func (dbPool *DbPool) FindAll(st string) []map[string]interface{} {
	rows, err := dbPool.db.Query(st)
	defer rows.Close()
	checkErr(err)

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	citems := make([]map[string]interface{}, 0, 10)
	for rows.Next() {
		record := make(map[string]interface{})
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		citems = append(citems, record)
	}
	return citems
}

func (dbPool *DbPool) FindOne(st string) map[string]interface{} {
	rows, err := dbPool.db.Query(st)
	defer rows.Close()
	checkErr(err)
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	if rows.Next() {
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
	}
	return record
}

func (dbPool *DbPool) Counts(st string) int {
	var cnt int
	_ = dbPool.db.QueryRow(st).Scan(&cnt)
	return cnt
}

func (dbPool *DbPool) MultiInsert(param []map[string]interface{}, tablename string) int64 {
	var keys []string
	var vals = []interface{}{}
	if len(param) > 0 {
		for key, _ := range param[0] {
			keys = append(keys, key)
		}
		fileds := "`" + strings.Join(keys, "`,`") + "`"
		sqlStr := fmt.Sprintf("REPLACE INTO %v (%v) VALUES ", tablename, fileds)

		for _, row := range param {
			sqlStr += "("
			for _, v := range keys {
				sqlStr += "?,"
				value := row[v]
				if value != nil {
					switch value.(type) {
					case int:
						vals = append(vals, strconv.Itoa(value.(int)))
					case int32, int64:
						vals = append(vals, strconv.FormatInt(value.(int64), 10))
					case string:
						vals = append(vals, EscapeString(value.(string)))
					case float32, float64:
						vals = append(vals, strconv.FormatFloat(value.(float64), 'f', -1, 64))
					default:
						vals = append(vals, "")
						fmt.Println("Replace into not type")
					}
				} else {
					vals = append(vals, "")
				}
				//vals = append(vals, row[v])
			}
			sqlStr = strings.TrimSuffix(sqlStr, ",")
			sqlStr += "),"
		}
		sqlStr = strings.TrimSuffix(sqlStr, ",")
		stmt, err := dbPool.db.Prepare(sqlStr)
		checkErr(err)
		result, err := stmt.Exec(vals...)
		if checkErr(err) {
			return -1
		}
		stmt.Close()
		affectLines, err := result.RowsAffected()
		return affectLines
	}
	return 0
}

func (dbPool *DbPool) Insert(param map[string]interface{}, tablename string) int64 {
	var keys []string
	var values []string
	for key, value := range param {
		keys = append(keys, key)
		if value != nil {
			switch value.(type) {
			case int:
				values = append(values, strconv.Itoa(value.(int)))
			case int32, int64:
				values = append(values, strconv.FormatInt(value.(int64), 10))
			case string:
				values = append(values, EscapeString(value.(string)))
			case float32, float64:
				values = append(values, strconv.FormatFloat(value.(float64), 'f', -1, 64))
			}
		} else {
			values = append(values, "")
		}
	}
	fileValue := "'" + strings.Join(values, "','") + "'"
	fileds := "`" + strings.Join(keys, "`,`") + "`"
	sql := fmt.Sprintf("REPLACE INTO %v (%v) VALUES (%v)", tablename, fileds, fileValue)
	result, err := dbPool.db.Exec(sql)
	if checkErr(err) {
		return -1
	}
	lastId, err := result.LastInsertId()
	checkErr(err)
	return lastId
}

func checkErr(err error) bool {
	if err != nil {
		if strings.Index(err.Error(), "Deadlock found when trying to get lock; try restarting transaction") > -1 {
			fmt.Println("==============Try Again======================")
			fmt.Println(err)
			time.Sleep(20000 * time.Millisecond)
			return true
		}
		panic(err)
	}
	return false
}
