package abstraction

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) CheckTrx(ctx *Context) *gorm.DB {
	if ctx.Trx != nil {
		return ctx.Trx.Db.WithContext(ctx.Request().Context()).Clauses(dbresolver.Write)
	}
	return r.Db.WithContext(ctx.Request().Context()).Clauses(dbresolver.Write)
}

func (r *Repository) Filter(ctx *Context, query *gorm.DB, payload interface{}) *gorm.DB {
	mVal := reflect.ValueOf(payload)
	mType := reflect.TypeOf(payload)

	for i := 0; i < mVal.NumField(); i++ {
		mValChild := mVal.Field(i)
		mTypeChild := mType.Field(i)

		for j := 0; j < mValChild.NumField(); j++ {
			val := mValChild.Field(j)

			if !val.IsNil() {
				if val.Kind() == reflect.Ptr {
					val = mValChild.Field(j).Elem()
				}

				key := mTypeChild.Type.Field(j).Tag.Get("query")
				filter := mTypeChild.Type.Field(j).Tag.Get("filter")

				switch filter {
				case "LIKE":
					query = query.Where(fmt.Sprintf("%s LIKE ?", key), "%"+val.String()+"%")
				case "ILIKE":
					query = query.Where(fmt.Sprintf("%s ILIKE ?", key), "%"+val.String()+"%")
				case "DATE":
					tmpDate, err := time.Parse("2006-01-02", val.String())
					if err != nil {
						continue
					}
					tmpStr := tmpDate.Format("2006-01-02")
					query = query.Where(fmt.Sprintf("DATE(%s) = ?", key), tmpStr)
				case "DATESTRING":
					query = query.Where(fmt.Sprintf("DATE(%s) = ?", key), val.String())
				case "IN":
					datas := strings.Split(strings.TrimSpace(val.String()), ",")
					query = query.Where(fmt.Sprintf("%s IN (?)", key), datas)
				case "CUSTOM":
					continue
				default:
					query = query.Where(fmt.Sprintf("%s = ?", key), val.Interface())
				}
			}
		}
	}

	return query
}

func (r *Repository) FilterWithTableName(ctx *Context, query *gorm.DB, payload interface{}, tableName string) *gorm.DB {
	mVal := reflect.ValueOf(payload)
	mType := reflect.TypeOf(payload)

	for i := 0; i < mVal.NumField(); i++ {
		mValChild := mVal.Field(i)
		mTypeChild := mType.Field(i)

		for j := 0; j < mValChild.NumField(); j++ {
			val := mValChild.Field(j)

			if !val.IsNil() {
				if val.Kind() == reflect.Ptr {
					val = mValChild.Field(j).Elem()
				}

				key := mTypeChild.Type.Field(j).Tag.Get("query")
				filter := mTypeChild.Type.Field(j).Tag.Get("filter")

				switch filter {
				case "LIKE":
					query = query.Where(fmt.Sprintf("%s.%s LIKE ?", tableName, key), "%"+val.String()+"%")
				case "ILIKE":
					query = query.Where(fmt.Sprintf("%s.%s ILIKE ?", tableName, key), "%"+val.String()+"%")
				case "DATE":
					tmpDate, err := time.Parse("2006-01-02", val.String())
					if err != nil {
						continue
					}
					tmpStr := tmpDate.Format("2006-01-02")
					query = query.Where(fmt.Sprintf("DATE(%s.%s) = ?", tableName, key), tmpStr)
				case "DATESTRING":
					query = query.Where(fmt.Sprintf("DATE(%s.%s) = ?", tableName, key), val.String())
				case "IN":
					datas := strings.Split(strings.TrimSpace(val.String()), ",")
					query = query.Where(fmt.Sprintf("%s.%s IN (?)", tableName, key), datas)
				case "CUSTOM":
					continue
				default:
					query = query.Where(fmt.Sprintf("%s.%s = ?", tableName, key), val.Interface())
				}
			}
		}
	}

	return query
}
