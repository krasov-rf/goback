package utils

import (
	"maps"
	"math/rand"
	"reflect"
	"unsafe"

	"github.com/lib/pq"
)

type TypeCode int

const (
	CODE_LETTERS TypeCode = iota
	CODE_NUBERS
)

const NUM_ALPHABET = "1234567890"
const LETTER_ALPHABET = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// получить рандомизированный ключ/строку
func GetRandKey(a TypeCode, n int) string {
	alphabet := LETTER_ALPHABET
	if a == CODE_NUBERS {
		alphabet = NUM_ALPHABET
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

// переводит структуру в мапу для использования в sql
// запросах, используется тэг db и рефлексия
func StructToMap(opt any, prefix string) map[string]any {
	res := map[string]any{}
	if opt == nil {
		return res
	}
	v := reflect.TypeOf(opt)
	reflectValue := reflect.ValueOf(opt)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := range v.NumField() {
		tag := v.Field(i).Tag.Get("db")
		field_type := v.Field(i).Type.Kind()
		if tag != "" && tag != "-" {
			field := reflectValue.Field(i).Interface()
			if field_type == reflect.Struct {
				res[prefix+tag] = StructToMap(field, prefix)
			} else {
				res[prefix+tag] = field
			}
		} else if tag != "-" && field_type == reflect.Struct {
			field := reflectValue.Field(i).Interface()
			maps.Copy(res, StructToMap(field, prefix))
		}
	}
	return res
}

// универсальная конвертация, пользоваться с осторожностью
// в данный момент конвертит []int в pq.Int64Array и наоборот
func Convert[T any, V any](value T) V {
	return *(*V)(unsafe.Pointer(&value))
}

// это чаще все го применяем
func ConvertInt64Array(value []int) pq.Int64Array {
	return Convert[[]int, pq.Int64Array](value)
}

func ConvertStringArray(value []string) pq.StringArray {
	return Convert[[]string, pq.StringArray](value)
}

// вернуть салайс ключей мапы
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// вернуть слайс значений мапы
func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// получить первый элемент из мапы, если она не пустая
func FirstM[K comparable, V any](record map[K]V) (V, bool) {
	var result V
	if len(record) == 0 {
		return result, false
	}
	for _, v := range record {
		return v, true
	}
	return result, false
}

// получить первый элемент из слайса, если он не пустой
func FirstS[V any](record []V) (V, bool) {
	var result V
	if len(record) == 0 {
		return result, false
	}
	result = record[0]
	return result, true
}
