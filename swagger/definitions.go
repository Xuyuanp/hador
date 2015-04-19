/*
 * Copyright 2015 Xuyuan Pang
 * Author: Xuyuan Pang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package swagger

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Definitions is global model definitions
type Definitions map[string]Schema

// AddModelFrom adds model definitions from model
func (d Definitions) AddModelFrom(model interface{}) {
	d.addModel(reflect.TypeOf(model), "")
}

func (d Definitions) addModel(st reflect.Type, nameOverride string) *Schema {
	modelName := st.Name()
	if nameOverride != "" {
		modelName = nameOverride
	}

	if d.isPrimitiveType(modelName) {
		return nil
	}

	if _, ok := d[modelName]; ok {
		return nil
	}

	if st.Kind() == reflect.Slice || st.Kind() == reflect.Array {
		return d.addModel(st.Elem(), "")
	}

	if st.Kind() != reflect.Struct {
		return nil
	}

	schema := Schema{
		Required:   []string{},
		Properties: map[string]Items{},
	}

	d[modelName] = schema

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if field.Anonymous {
			nestedSchema, ok := d[field.Type.Name()]
			if !ok {
				nestedSchema = *d.addModel(field.Type, "")
			}
			d.mergeSchema(&schema, &nestedSchema)
			continue
		}
		jsonName, prop := d.buildProperty(field, &schema, modelName)
		if jsonName != "" {
			if d.isPropertyRequired(field) {
				schema.Required = append(schema.Required, jsonName)
			}
			schema.Properties[jsonName] = prop
		}
	}

	d[modelName] = schema

	return &schema
}

func (d Definitions) buildProperty(field reflect.StructField, schema *Schema, modelName string) (jsonName string, prop Items) {
	jsonName = d.jsonNameOfField(field)
	if len(jsonName) == 0 {
		return
	}
	fieldType := field.Type

	marshalerType := reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	if fieldType.Implements(marshalerType) {
		prop.Type = "string"
		prop.Format = d.jsonSchemaFormat(fieldType.String())
		return
	}

	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if len(s) > 1 && s[1] == "string" {
			prop.Type = "string"
			return
		}
	}

	if d.isPrimitiveType(fieldType.String()) {
		prop.Type = d.jsonSchemaType(fieldType.String())
		prop.Format = d.jsonSchemaFormat(fieldType.String())
		return
	}

	switch fieldType.Kind() {
	case reflect.Struct:
		return d.buildStructProperty(field, jsonName, modelName)
	case reflect.Slice, reflect.Array:
	case reflect.Ptr:
	case reflect.Map:
		prop.Type = "any"
		return
	}
	return
}

func (d Definitions) buildPointerProperty(field reflect.StructField, jsonName, modelName string) (pName string, prop Items) {

	return
}

func (d Definitions) buildStructProperty(field reflect.StructField, jsonName, modelName string) (sName string, prop Items) {
	fieldType := field.Type
	d.addModel(fieldType, "")
	prop.Ref = "#/definitions/" + fieldType.Name()
	sName = jsonName
	return
}

func (d Definitions) buildArrayProperty(field reflect.StructField, jsonName, modelName string) (ajName string, prop Items) {
	return
}

func (d Definitions) isPrimitiveType(modelName string) bool {
	return strings.Contains("uint8 uint16 uint32 uint64 int int8 int16 int32 int64 float32 float64 bool string byte rune time.Time", modelName)
}

func (d Definitions) jsonNameOfField(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if s[0] == "-" {
			return ""
		}
		if s[0] != "" {
			return s[0]
		}
	}
	return field.Name
}

func (d Definitions) jsonSchemaType(modelName string) string {
	schemaMap := map[string]string{
		"uint8":     "integer",
		"uint16":    "integer",
		"uint32":    "integer",
		"uint64":    "integer",
		"int":       "integer",
		"int8":      "integer",
		"int16":     "integer",
		"int32":     "integer",
		"int64":     "integer",
		"byte":      "integer",
		"float32":   "number",
		"float64":   "number",
		"bool":      "boolean",
		"time.Time": "string",
	}
	if t, ok := schemaMap[modelName]; ok {
		return t
	}
	return modelName
}

func (d Definitions) jsonSchemaFormat(modelName string) string {
	schemaMap := map[string]string{
		"int":       "int32",
		"int32":     "int32",
		"int64":     "int64",
		"byte":      "byte",
		"uint8":     "byte",
		"float64":   "double",
		"float32":   "float",
		"time.Time": "date-time",
	}
	if f, ok := schemaMap[modelName]; ok {
		return f
	}
	return ""
}

func (d Definitions) mergeSchema(a *Schema, b *Schema) {
	a.Required = append(a.Required, b.Required...)
	for t, i := range b.Properties {
		if _, ok := a.Properties[t]; !ok {
			a.Properties[t] = i
		}
	}
}

func (d Definitions) isPropertyRequired(field reflect.StructField) bool {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if len(s) > 1 && s[1] == "omitempty" {
			return false
		}
	}
	return true
}
