package internal

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func stringify(rv reflect.Value) string {
	switch rv.Kind() {
	case reflect.Map:
		value := rv.Interface().(map[string]string)
		stringValue := []string{}
		for k, v := range value {
			stringValue = append(stringValue, fmt.Sprintf("%s=%s", k, v))
		}
		return strings.Join(stringValue, "\n")
	case reflect.Slice:
		stringValue := rv.Interface().([]string)
		return strings.Join(stringValue, ",")
	case reflect.Bool:
		boolValue := rv.Interface().(bool)
		if boolValue {
			return "true"
		} else {
			return "false"
		}
	default:
		return rv.String()
	}
}

func inspect(t interface{}) map[string]string {
	s := reflect.ValueOf(t).Elem()
	typeOfT := s.Type()
	m := make(map[string]string)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if name, ok := typeOfT.Field(i).Tag.Lookup("view"); ok {
			k := strings.Split(name, ",")[0]
			m[k] = stringify(f)
		} else {
			m[typeOfT.Field(i).Name] = stringify(f)
		}
	}
	return m
}

func sortKeysByTag(t reflect.Type, keys []string) {
	typeOfT := t
	order := map[string]int{}
	for i := 0; i < t.NumField(); i++ {
		if name, ok := typeOfT.Field(i).Tag.Lookup("view"); ok {
			tags := strings.Split(name, ",")
			if len(tags) > 1 {
				if j, err := strconv.Atoi(tags[1]); err == nil {
					order[tags[0]] = j
				}
			}
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return order[keys[i]] < order[keys[j]]
	})
}

func PrettyMapAsArray(o map[string]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	for k, v := range o {
		table.Append([]string{k, v})
	}
	table.Render() // Send output
}

func PrettyArray(o []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key"})
	for _, v := range o {
		table.Append([]string{v})
	}
	table.Render() // Send output
}

func PrettyStruct(o interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	p := inspect(o)
	keys := []string{}
	for k := range p {
		keys = append(keys, k)
	}
	sortKeysByTag(reflect.TypeOf(o).Elem(), keys)
	for _, k := range keys {
		table.Append([]string{fmt.Sprintf("%s:", k), p[k]})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.Render()              // Send output
}

func PrettyStructArray(slice interface{}, class reflect.Type) {

	if reflect.TypeOf(slice).Elem().Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	s := reflect.ValueOf(slice).Elem()
	// to map array
	ret := make([]map[string]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = inspect(s.Index(i).Addr().Interface())
	}

	keys := []string{}
	for i := 0; i < class.NumField(); i++ {
		if name, ok := class.Field(i).Tag.Lookup("view"); ok {
			tags := strings.Split(name, ",")
			if len(tags) > 1 {
				if _, err := strconv.Atoi(tags[1]); err == nil {
					keys = append(keys, tags[0])
				}
			}
		}
	}
	sortKeysByTag(class, keys)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(keys)
	for _, obj := range ret {
		values := []string{}
		for _, key := range keys {
			values = append(values, obj[key])
		}
		table.Append(values)
	}
	table.Render() // Send output
}
