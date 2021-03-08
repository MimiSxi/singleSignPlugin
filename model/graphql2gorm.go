package model

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func defaultNamer(name string) string {
	const (
		lower = false
		upper = true
	)

	var (
		value                                    = name
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	return s
}

var op2string map[string]string = map[string]string{
	"_eq":   "=",
	"_neq":  "<>",
	"_gt":   ">",
	"_lt":   "<",
	"_gte":  ">=",
	"_lte":  "<=",
	"_like": "like",
	"_in":   "in",
	"_nin":  "not in",
}

func GenSelet(db *gorm.DB, params graphql.ResolveParams) *gorm.DB {
	if params.Args["first"] == nil {
		db = db.Limit(10)
	} else {
		i := params.Args["first"].(int)
		db = db.Limit(i)
	}

	if params.Args["skip"] != nil {
		i := params.Args["skip"].(int)
		db = db.Offset(i)
	}

	db = GenWhere(db, params)
	if params.Args["orderby"] != nil {
		if orderby, ok := params.Args["orderby"].([]interface{}); ok {
			for _, v := range orderby {
				if item, ok := v.(map[string]interface{}); ok {
					for k1, v1 := range item {
						db = db.Order(fmt.Sprintf("%s %s", defaultNamer(k1), v1))
					}
				}
			}
		}

	}
	return db
}

func GenWhere(db *gorm.DB, params graphql.ResolveParams) *gorm.DB {
	if params.Args["where"] != nil {
		where := params.Args["where"].(map[string]interface{})
		for k, v := range where {
			k = defaultNamer(k)
			if d, ok := v.(map[string]interface{}); ok {
				for k1, v1 := range d {
					if op, ok := op2string[k1]; !ok {
						logrus.Warnf("未支持 %s", k1)
					} else {
						switch op {
						case "in", "nin":
							db = db.Where(fmt.Sprintf("%s %s (?)", k, op), v1)
						case "like":
							db = db.Where(fmt.Sprintf("%s %s ?", k, op), v1)
						default:
							db = db.Where(fmt.Sprintf("%s %s ?", k, op), v1)
						}
					}
				}
			}
		}
	}
	return db
}
