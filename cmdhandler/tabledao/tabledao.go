package tabledao

import (
	"cmdtools/domain/repository"
	"log"
	"regexp"
)

const (
	// 判断类型
	int64Type   = `int`
	float64Type = `(decimal)|(double)|(float)|(numeric)`
)

type TableDesc struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

func DescTable(table string) []TableDesc {
	var fields []TableDesc
	res := repository.GormDB.Raw("DESC " + table).Scan(&fields)
	if res.Error != nil {
		log.Fatalln("DESC "+table+" error:", res.Error)
	}
	convertToGolangType(&fields)
	return fields
}

func convertToGolangType(fields *[]TableDesc) {
	int64TypeRegx := regexp.MustCompile(int64Type)
	float64Type := regexp.MustCompile(float64Type)
	for index := range *fields {
		switch {
		case int64TypeRegx.MatchString((*fields)[index].Type):
			(*fields)[index].Type = "int64"
		case float64Type.MatchString((*fields)[index].Type):
			(*fields)[index].Type = "float64"
		default:
			(*fields)[index].Type = "string"
		}
	}
}
