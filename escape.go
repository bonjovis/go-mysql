/*
*
* Author: Hui Ye - <bonjovis@163.com>
*
* Last modified: 2017-03-17 09:21
*
* Filename: escape.go
*
* Copyright (c) 2016 JOVI
*
 */

package mysql

func EscapeString(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]

		escape = 0

		switch c {
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}
