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
	if st.Kind() == reflect.Ptr || st.Kind() == reflect.Array || st.Kind() == reflect.Slice {
		return d.addModel(st.Elem(), nameOverride)
	}

	modelName := st.String()
	if nameOverride != "" {
		modelName = nameOverride
	}
	if modelName == "" {
		return nil
	}

	if schema, ok := d[modelName]; ok {
		return &schema
	}

	if st.Kind() != reflect.Struct {
		return nil
	}

	schema := Schema{
		Properties: map[string]Items{},
		Required:   []string{},
	}

	d[modelName] = schema

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if field.Anonymous {
			nestedSchema, ok := d[field.Type.String()]
			if !ok {
				nestedSchema = *d.addModel(field.Type, "")
			}
			d.mergeSchema(&schema, &nestedSchema)
			continue
		}
		jsonName, prop := d.buildProperty(field, "")
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

func (d Definitions) buildProperty(field reflect.StructField, nameOverride string) (jsonName string, prop Items) {
	jsonName = d.jsonNameOfField(field)
	if nameOverride != "" {
		jsonName = nameOverride
	}
	if jsonName == "" {
		return
	}

	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if len(s) > 1 && s[1] == "string" {
			prop.Type = "string"
			return
		}
	}

	prop = d.buildFieldType(field.Type)
	return
}

func (d Definitions) buildFieldType(st reflect.Type) (prop Items) {
	if d.isPrimitiveType(st.String()) {
		prop.Type = d.jsonSchemaType(st.String())
		prop.Format = d.jsonSchemaFormat(st.String())
		return
	}
	switch st.Kind() {
	case reflect.Struct:
		d.addModel(st, "")
		prop.Ref = "#/definitions/" + st.String()
	case reflect.Slice, reflect.Array:
		prop.Type = "array"
		itemsprop := d.buildFieldType(st.Elem())
		prop.Items = &itemsprop
	case reflect.Ptr:
		return d.buildFieldType(st.Elem())
	}
	return
}

func (d Definitions) isPrimitiveType(modelName string) bool {
	return modelName != "" && strings.Contains("uint8 uint16 uint32 uint64 int int8 int16 int32 int64 float32 float64 bool string byte rune time.Time", modelName)
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
