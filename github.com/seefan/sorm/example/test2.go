package main

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// DefaultTagName 默认 tag_name
	DefaultTagName = "json"
)

// Person 代表 person
type Person struct {
	Age    int    `json:"age"`
	Name   string `json:"name"`
	Offer  string `json:"-"`
	Gender string
}

// "age,omitempty"
func parseTags(s string) (string, []string) {
	if len(s) == 0 {
		return "", nil
	}
	sl := strings.Split(s, ",")
	return sl[0], sl[1:]
}

// ConvertStructToMap struct to map json tag
func ConvertStructToMap(val interface{}) (map[string]interface{}, error) {
	var (
		rv     = reflect.ValueOf(val)
		rvt    = rv.Type()
		fields []reflect.StructField
		result = map[string]interface{}{}
	)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not struct")
	}
	// 2. get all struct_field
	for i := 0; i < rvt.NumField(); i++ {
		field := rvt.Field(i)
		// unexported field
		if field.PkgPath != "" {
			continue
		}
		// ignored tag
		if tag := field.Tag.Get(DefaultTagName); tag == "-" {
			continue
		}
		fields = append(fields, field)
	}
	// 3. fill map
	for _, field := range fields {
		fieldName := field.Name
		fieldValue := rv.FieldByName(fieldName)
		// tagName, tagOpts
		tagName, _ := parseTags(field.Tag.Get(DefaultTagName))
		// TODO: add some json tag option process
		if len(tagName) == 0 {
			tagName = fieldName
		}
		result[tagName] = fieldValue.Interface()
	}
	return result, nil
}

// SetField struct field set value
func SetField(sval interface{}, name string, val interface{}) error {
	var (
		sv = reflect.ValueOf(sval)
		fv reflect.Value
		ft reflect.Type
	)
	for sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	// TODO: tag --> field_name
	// TODO: sv.FieldByNameFunc
	fv = sv.FieldByName(name)
	if !fv.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !fv.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	ft = fv.Type()
	if ft != reflect.ValueOf(val).Type() {
		return fmt.Errorf("Provided value type didn't match obj field type")
	}
	fv.Set(reflect.ValueOf(val))
	return nil
}

// ConvertMapToStruct params map fill struct
func ConvertMapToStruct(params map[string]interface{}, val interface{}) {
	for field, fieldVal := range params {
		if err := SetField(val, field, fieldVal); err != nil {
			fmt.Println(err)
		}
	}
}
func main() {
	var p = Person{
		Age:    18,
		Name:   "wangxiyang",
		Offer:  "didi",
		Gender: "nan",
	}
	mm, _ := ConvertStructToMap(p)
	fmt.Println(mm)

	var p2 Person
	var params = map[string]interface{}{
		"Age":    18,
		"Name":   "wangxiyang",
		"Offer":  "didi",
		"Gender": "nan",
	}
	// 必须传指针方能 can_set
	ConvertMapToStruct(params, &p2)
	fmt.Printf("%+v\n", p2)
}
