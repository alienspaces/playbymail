package tag

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/shopspring/decimal"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

var nonRecurseTypeNames = set.New(
	getTypeName(sql.NullString{}),
	getTypeName(sql.NullInt32{}),
	getTypeName(sql.NullInt64{}),
	getTypeName(sql.NullFloat64{}),
	getTypeName(sql.NullTime{}),
	getTypeName(sql.NullBool{}),
	getTypeName(time.Time{}),
	getTypeName(decimal.Decimal{}),
	getTypeName(decimal.NullDecimal{}),
)

func getTypeName(i any) string {
	return reflect.TypeOf(i).Name()
}

func GetFieldTagValues(entity any, tag string, fp ...ExcludeFilterPredicate) []string {
	t := reflect.TypeOf(entity)

	values := &[]string{}
	getValuesRecur(t, tag, values, fp...)
	return *values
}

func getValuesRecur(t reflect.Type, tag string, values *[]string, fp ...ExcludeFilterPredicate) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if shouldRecurse(field) {
			getValuesRecur(field.Type, tag, values, fp...)
		} else if shouldExclude(field, fp...) {
			continue
		} else if attr, ok := shouldGetTagValue(field, tag); ok && attr != "-" {
			*values = append(*values, attr)
		}
	}
}

func shouldRecurse(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Struct && isRecurseFieldType(field.Type.Name())
}

func isRecurseFieldType(name string) bool {
	return !nonRecurseTypeNames.Has(name)
}

func shouldGetTagValue(field reflect.StructField, tag string) (string, bool) {
	return field.Tag.Lookup(tag)
}
