package server

import (
	"os"
	"reflect"
	"slices"
	"strings"
)

type StructFieldInfo struct {
	GoName     string
	TSName     string
	JSONName   string
	OmitEmpty  bool
	Ignored    bool
	Exported   bool
	Anonymous  bool
	Type       reflect.Type
	TagRawJSON string
}

func (s *Server) exportInterfacesToTS() error {
	allInterfaces := ""

	exportedTypes := []reflect.Type{}

	for _, paths := range s.Paths {
		for _, p := range paths {
			for _, t := range p.Info.ExportTypes {
				if !slices.Contains(exportedTypes, t) {
					exportedTypes = append(exportedTypes, t)
					allInterfaces += createTSInterface(t)
				}
			}
		}
	}

	if allInterfaces != "" {
		os.WriteFile(s.ExportTypesLocation, []byte(allInterfaces), os.ModePerm)
	}
	return nil
}

func createTSInterface(t reflect.Type) string {
	infos := parseStructFields(t)

	var b strings.Builder
	b.WriteString("export interface ")
	b.WriteString(t.Name())
	b.WriteString(" {\n")

	for _, i := range infos {
		if i.Ignored {
			continue
		}

		tsType, optFromPtr := goTypeToTsType(i.Type)
		optional := optFromPtr || i.OmitEmpty
		prop := i.TSName

		b.WriteString("  ")
		b.WriteString(prop)
		if optional {
			b.WriteString("?: ")
		} else {
			b.WriteString(": ")
		}
		b.WriteString(tsType)
		b.WriteString(";\n")
	}

	b.WriteString("}\n")

	return b.String()
}

func goTypeToTsType(t reflect.Type) (ts string, optional bool) {
	optional = false

	if t.Kind() == reflect.Ptr {
		optional = true
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Bool:
		return "boolean", optional
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number", optional
	case reflect.String:
		return "string", optional
	case reflect.Slice, reflect.Array:
		elem, _ := goTypeToTsType(t.Elem())
		return elem + "[]", optional
	case reflect.Map:
		// TS index signature only works well with string/number keys
		key, _ := goTypeToTsType(t.Key())
		val, _ := goTypeToTsType(t.Elem())
		if key != "string" && key != "number" {
			key = "string" // pragmatic fallback
		}
		return "{ [key: " + key + "]: " + val + " }", optional
	case reflect.Struct:
		// special-case common structs
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return "string", optional // or "Date"
		}
		// otherwise reference another interface by name
		return t.Name(), optional
	case reflect.Interface:
		return "any", optional
	default:
		return "any", optional
	}
}

func parseJSONTag(tag string) (name string, omitEmpty bool, ignored bool) {
	if tag == "" {
		return "", false, false
	}
	if tag == "-" {
		return "", false, true
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, p := range parts[1:] {
		if p == "omitempty" {
			omitEmpty = true
		}
	}
	return name, omitEmpty, false
}

func parseStructFields(t reflect.Type) []StructFieldInfo {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	var out []StructFieldInfo
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		// skip unexported (PkgPath != "" means unexported, except embedded sometimes)
		exported := sf.PkgPath == ""
		if !exported && !sf.Anonymous {
			continue
		}

		raw := sf.Tag.Get("json")
		name, omitempty, ignored := parseJSONTag(raw)

		jsonName := name
		if jsonName == "" && !ignored {
			jsonName = sf.Name
		}

		fi := StructFieldInfo{
			GoName:     sf.Name,
			TSName:     jsonName,
			JSONName:   jsonName,
			OmitEmpty:  omitempty,
			Ignored:    ignored,
			Exported:   exported,
			Anonymous:  sf.Anonymous,
			Type:       sf.Type,
			TagRawJSON: raw,
		}
		out = append(out, fi)
	}
	return out
}
