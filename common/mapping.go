package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func CamelCase(name string, firstUpper bool) string {
	array := []byte(name)
	if len(array) == 0 {
		return ""
	}
	rspArray := make([]byte, len(array))
	if firstUpper {
		copy(rspArray[:1], strings.ToUpper(string(array[:1])))
	} else {
		copy(rspArray[:1], strings.ToLower(string(array[:1])))
	}
	copy(rspArray[1:], array[1:])
	return string(rspArray)
}

func ObjectToMap(o interface{}) map[string]interface{} {
	if o == nil {
		return nil
	}
	value := reflect.ValueOf(o)
	if value.Kind() != reflect.Ptr {
		return nil
	}
	elem := value.Elem()
	relType := elem.Type()
	m := make(map[string]interface{})
	for i := 0; i < relType.NumField(); i++ {
		field := relType.Field(i)
		if elem.Field(i).IsZero() {
			continue
		}
		m[CamelCase(field.Name, false)] = elem.Field(i).Interface()
	}
	return m
}

// AssertString convert v to string value
func AssertString(v interface{}) string {
	if v == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}

// ValidatePtr validate v is a ptr value
func ValidatePtr(v *reflect.Value) error {
	// sequence is very important, IsNil must be called after checking Kind() with reflect.Ptr,
	// panic otherwise
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("not a valid pointer: %v", v)
	}

	return nil
}

func LoadCustomFieldToMap(src interface{}, fields ...string) map[string]interface{} {
	rsp := LoadCustomField(src, fields...)
	if rsp == nil {
		return map[string]interface{}{}
	}
	return rsp.(map[string]interface{})
}

func LoadCustomField(src interface{}, fields ...string) interface{} {
	typeSrc := reflect.TypeOf(src)
	valueSrc := reflect.ValueOf(src)

	if v, ok := src.(reflect.Value); ok {
		valueSrc = v
		typeSrc = v.Type()
	}
	if typeSrc.Kind() == reflect.Ptr {
		valueSrc = valueSrc.Elem()
	}
	k := valueSrc.Kind()
	switch k {
	case reflect.Array, reflect.Slice:
		len := valueSrc.Len()
		retSliceMap := make([]map[string]interface{}, 0)
		if len == 0 {
			return retSliceMap
		}
		for i := 0; i < len; i++ {
			v := valueSrc.Index(i)
			retSliceMap = append(retSliceMap, (LoadCustomField(v, fields...)).(map[string]interface{}))
		}
		return retSliceMap
	case reflect.Struct:
		retSliceMap := make(map[string]interface{})
		for _, filed := range fields {
			f := valueSrc.FieldByName(filed)
			if !f.IsValid() {
				continue
			}
			v := f.Interface()
			if t, ok := v.(time.Time); ok {
				v = t.Local().Format("2006-01-02 15:04:05")
			}
			retSliceMap[CamelCase(filed, false)] = v
		}
		return retSliceMap
	default:
		return src
	}
	return src
}

func AppendCustomField(src interface{}, options map[string]interface{}) interface{} {
	var mapSrc map[string]interface{}
	var ok bool
	mapSrc, ok = src.(map[string]interface{})
	if !ok {
		JsonUnmarshal(JsonAssertString(src), &mapSrc)
	}
	for field, value := range options {
		mapSrc[CamelCase(field, false)] = value
	}
	return mapSrc
}

/*

json 格式化

*/

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return formatError(string(data), err)
	}

	return nil
}

func UnmarshalFromString(str string, v interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(str))
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return formatError(str, err)
	}

	return nil
}

func UnmarshalFromReader(reader io.Reader, v interface{}) error {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := json.NewDecoder(teeReader)
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return formatError(buf.String(), err)
	}

	return nil
}

func unmarshalUseNumber(decoder *json.Decoder, v interface{}) error {
	decoder.UseNumber()
	return decoder.Decode(v)
}

func formatError(v string, err error) error {
	return fmt.Errorf("string: `%s`, error: `%s`", v, err.Error())
}

type ReflectVal struct {
	T reflect.Type
	V reflect.Value
}

/*
	拷贝当前对象到目标对象，具有相同属性的值
*/
func CopyObject(src, dst interface{}) {
	var srcMap = make(map[string]ReflectVal)

	vs := reflect.ValueOf(src)
	ts := reflect.TypeOf(src)
	vd := reflect.ValueOf(dst)
	td := reflect.TypeOf(dst)

	ls := vs.Elem().NumField()
	for i := 0; i < ls; i++ {
		srcMap[ts.Elem().Field(i).Name] = ReflectVal{
			T: vs.Elem().Field(i).Type(),
			V: vs.Elem().Field(i),
		}
	}

	ld := vd.Elem().NumField()
	for i := 0; i < ld; i++ {
		n := td.Elem().Field(i).Name
		t := vd.Elem().Field(i).Type()
		if v, ok := srcMap[n]; ok && v.T == t && vd.Elem().Field(i).CanSet() {
			vd.Elem().Field(i).Set(v.V)
		}
	}
}
