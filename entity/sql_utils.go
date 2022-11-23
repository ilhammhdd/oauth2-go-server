package entity

import (
	"encoding/json"
	"strings"
)

type DBSet []byte

func (dbs DBSet) MarshalJSON() ([]byte, error) {
	dbsStr := string(dbs)
	dbsStrSplit := strings.Split(dbsStr, ",")
	return json.Marshal(dbsStrSplit)
}

func (dbs *DBSet) UnmarshalJSON(b []byte) error {
	var result DBSet
	for idx := range b {
		if b[idx] != '[' && b[idx] != ']' && b[idx] != '"' && b[idx] != '\n' && b[idx] != '\t' {
			result = append(result, b[idx])
		}
	}
	*dbs = result
	return nil
}
