# go-mysql lib
Golang language integration of mysql operations

## Table of contents:
- [Get Started](#get-started)
- [Examples](#examples)


### Get Started
#### Installation

```sh
$ go get github.com/go-sql-driver/mysql
$ go get github.com/bonjovis/go-mysql
```


#### Examples
```go
import (
	"fmt"
	"github.com/bonjovis/go-mysql"
)

func main() {
  dbHost := "user:password@tcp(dbhost:3306)/dbname"
  maxOpenConns := 200
  maxIdleConns := 100
  dbPool := mysql.NewDatabaseConnectionPool(dbHost, maxOpenConns, maxIdleConns)
  
  //query
  list := dbPool.FindAll("select * from list")
  //counts
  counts := dbPool.Counts("select count(1) from list")
  //update
  var vals = []interface{}{}
  dbPool.Update("update test set abc=1",vals)
  //insert
  vo := make(map[string]interface{})
  tableName := "test"
  vo["id"] = 1
  vo["name"] = "test"
  ret := dbPool.Insert(vo, tableName)
  //MultiInsert
  var list = []map[string]interface{}
  list = append(list, vo)
  ret = dbPool.MultiInsert(list, tableName)
}
```


### License
MIT
