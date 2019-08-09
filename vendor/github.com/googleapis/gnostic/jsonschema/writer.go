<<<<<<< HEAD
// Copyright 2017 Google LLC. All Rights Reserved.
=======
// Copyright 2017 Google Inc. All Rights Reserved.
>>>>>>> 79bfea2d (update vendor)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonschema

import (
	"fmt"
<<<<<<< HEAD

	"gopkg.in/yaml.v3"
=======
	"gopkg.in/yaml.v2"
>>>>>>> 79bfea2d (update vendor)
)

const indentation = "  "

<<<<<<< HEAD
func renderMappingNode(node *yaml.Node, indent string) (result string) {
	result = "{\n"
	innerIndent := indent + indentation
	for i := 0; i < len(node.Content); i += 2 {
		// first print the key
		key := node.Content[i].Value
		result += fmt.Sprintf("%s\"%+v\": ", innerIndent, key)
		// then the value
		value := node.Content[i+1]
		switch value.Kind {
		case yaml.ScalarNode:
			result += "\"" + value.Value + "\""
		case yaml.MappingNode:
			result += renderMappingNode(value, innerIndent)
		case yaml.SequenceNode:
			result += renderSequenceNode(value, innerIndent)
		default:
			result += fmt.Sprintf("???MapItem(Key:%+v, Value:%T)", value, value)
		}
		if i < len(node.Content)-2 {
			result += ","
		}
		result += "\n"
=======
func renderMap(info interface{}, indent string) (result string) {
	result = "{\n"
	innerIndent := indent + indentation
	switch pairs := info.(type) {
	case yaml.MapSlice:
		for i, pair := range pairs {
			// first print the key
			result += fmt.Sprintf("%s\"%+v\": ", innerIndent, pair.Key)
			// then the value
			switch value := pair.Value.(type) {
			case string:
				result += "\"" + value + "\""
			case bool:
				if value {
					result += "true"
				} else {
					result += "false"
				}
			case []interface{}:
				result += renderArray(value, innerIndent)
			case yaml.MapSlice:
				result += renderMap(value, innerIndent)
			case int:
				result += fmt.Sprintf("%d", value)
			case int64:
				result += fmt.Sprintf("%d", value)
			case []string:
				result += renderStringArray(value, innerIndent)
			default:
				result += fmt.Sprintf("???MapItem(Key:%+v, Value:%T)", value, value)
			}
			if i < len(pairs)-1 {
				result += ","
			}
			result += "\n"
		}
	default:
		// t is some other type that we didn't name.
>>>>>>> 79bfea2d (update vendor)
	}

	result += indent + "}"
	return result
}

<<<<<<< HEAD
func renderSequenceNode(node *yaml.Node, indent string) (result string) {
	result = "[\n"
	innerIndent := indent + indentation
	for i := 0; i < len(node.Content); i++ {
		item := node.Content[i]
		switch item.Kind {
		case yaml.ScalarNode:
			result += innerIndent + "\"" + item.Value + "\""
		case yaml.MappingNode:
			result += innerIndent + renderMappingNode(item, innerIndent) + ""
		default:
			result += innerIndent + fmt.Sprintf("???ArrayItem(%+v)", item)
		}
		if i < len(node.Content)-1 {
=======
func renderArray(array []interface{}, indent string) (result string) {
	result = "[\n"
	innerIndent := indent + indentation
	for i, item := range array {
		switch item := item.(type) {
		case string:
			result += innerIndent + "\"" + item + "\""
		case bool:
			if item {
				result += innerIndent + "true"
			} else {
				result += innerIndent + "false"
			}
		case yaml.MapSlice:
			result += innerIndent + renderMap(item, innerIndent) + ""
		default:
			result += innerIndent + fmt.Sprintf("???ArrayItem(%+v)", item)
		}
		if i < len(array)-1 {
>>>>>>> 79bfea2d (update vendor)
			result += ","
		}
		result += "\n"
	}
	result += indent + "]"
	return result
}

func renderStringArray(array []string, indent string) (result string) {
	result = "[\n"
	innerIndent := indent + indentation
	for i, item := range array {
		result += innerIndent + "\"" + item + "\""
		if i < len(array)-1 {
			result += ","
		}
		result += "\n"
	}
	result += indent + "]"
	return result
}

<<<<<<< HEAD
// Render renders a yaml.Node as JSON
func Render(node *yaml.Node) string {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 1 {
			return Render(node.Content[0])
		}
	} else if node.Kind == yaml.MappingNode {
		return renderMappingNode(node, "") + "\n"
	} else if node.Kind == yaml.SequenceNode {
		return renderSequenceNode(node, "") + "\n"
	}
	return ""
}

func (object *SchemaNumber) nodeValue() *yaml.Node {
	if object.Integer != nil {
		return nodeForInt64(*object.Integer)
	} else if object.Float != nil {
		return nodeForFloat64(*object.Float)
=======
func Render(info yaml.MapSlice) string {
	return renderMap(info, "") + "\n"
}

func (object *SchemaNumber) jsonValue() interface{} {
	if object.Integer != nil {
		return object.Integer
	} else if object.Float != nil {
		return object.Float
>>>>>>> 79bfea2d (update vendor)
	} else {
		return nil
	}
}

<<<<<<< HEAD
func (object *SchemaOrBoolean) nodeValue() *yaml.Node {
	if object.Schema != nil {
		return object.Schema.nodeValue()
	} else if object.Boolean != nil {
		return nodeForBoolean(*object.Boolean)
=======
func (object *SchemaOrBoolean) jsonValue() interface{} {
	if object.Schema != nil {
		return object.Schema.jsonValue()
	} else if object.Boolean != nil {
		return *object.Boolean
>>>>>>> 79bfea2d (update vendor)
	} else {
		return nil
	}
}

<<<<<<< HEAD
func nodeForStringArray(array []string) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, item := range array {
		content = append(content, nodeForString(item))
	}
	return nodeForSequence(content)
}

func nodeForSchemaArray(array []*Schema) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, item := range array {
		content = append(content, item.nodeValue())
	}
	return nodeForSequence(content)
}

func (object *StringOrStringArray) nodeValue() *yaml.Node {
	if object.String != nil {
		return nodeForString(*object.String)
	} else if object.StringArray != nil {
		return nodeForStringArray(*(object.StringArray))
=======
func (object *StringOrStringArray) jsonValue() interface{} {
	if object.String != nil {
		return *object.String
	} else if object.StringArray != nil {
		array := make([]interface{}, 0)
		for _, item := range *(object.StringArray) {
			array = append(array, item)
		}
		return array
>>>>>>> 79bfea2d (update vendor)
	} else {
		return nil
	}
}

<<<<<<< HEAD
func (object *SchemaOrStringArray) nodeValue() *yaml.Node {
	if object.Schema != nil {
		return object.Schema.nodeValue()
	} else if object.StringArray != nil {
		return nodeForStringArray(*(object.StringArray))
=======
func (object *SchemaOrStringArray) jsonValue() interface{} {
	if object.Schema != nil {
		return object.Schema.jsonValue()
	} else if object.StringArray != nil {
		array := make([]interface{}, 0)
		for _, item := range *(object.StringArray) {
			array = append(array, item)
		}
		return array
>>>>>>> 79bfea2d (update vendor)
	} else {
		return nil
	}
}

<<<<<<< HEAD
func (object *SchemaOrSchemaArray) nodeValue() *yaml.Node {
	if object.Schema != nil {
		return object.Schema.nodeValue()
	} else if object.SchemaArray != nil {
		return nodeForSchemaArray(*(object.SchemaArray))
=======
func (object *SchemaOrSchemaArray) jsonValue() interface{} {
	if object.Schema != nil {
		return object.Schema.jsonValue()
	} else if object.SchemaArray != nil {
		array := make([]interface{}, 0)
		for _, item := range *(object.SchemaArray) {
			array = append(array, item.jsonValue())
		}
		return array
>>>>>>> 79bfea2d (update vendor)
	} else {
		return nil
	}
}

<<<<<<< HEAD
func (object *SchemaEnumValue) nodeValue() *yaml.Node {
	if object.String != nil {
		return nodeForString(*object.String)
	} else if object.Bool != nil {
		return nodeForBoolean(*object.Bool)
=======
func (object *SchemaEnumValue) jsonValue() interface{} {
	if object.String != nil {
		return *object.String
	} else if object.Bool != nil {
		return *object.Bool
>>>>>>> 79bfea2d (update vendor)
	} else {
		return nil
	}
}

<<<<<<< HEAD
func nodeForNamedSchemaArray(array *[]*NamedSchema) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, pair := range *(array) {
		content = appendPair(content, pair.Name, pair.Value.nodeValue())
	}
	return nodeForMapping(content)
}

func nodeForNamedSchemaOrStringArray(array *[]*NamedSchemaOrStringArray) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, pair := range *(array) {
		content = appendPair(content, pair.Name, pair.Value.nodeValue())
	}
	return nodeForMapping(content)
}

func nodeForSchemaEnumArray(array *[]SchemaEnumValue) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, item := range *array {
		content = append(content, item.nodeValue())
	}
	return nodeForSequence(content)
}

func nodeForMapping(content []*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: content,
	}
}

func nodeForSequence(content []*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: content,
	}
}

func nodeForString(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}
}

func nodeForBoolean(value bool) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!bool",
		Value: fmt.Sprintf("%t", value),
	}
}

func nodeForInt64(value int64) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int",
		Value: fmt.Sprintf("%d", value),
	}
}

func nodeForFloat64(value float64) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!float",
		Value: fmt.Sprintf("%f", value),
	}
}

func appendPair(nodes []*yaml.Node, name string, value *yaml.Node) []*yaml.Node {
	nodes = append(nodes, nodeForString(name))
	nodes = append(nodes, value)
	return nodes
}

func (schema *Schema) nodeValue() *yaml.Node {
	n := &yaml.Node{Kind: yaml.MappingNode}
	content := make([]*yaml.Node, 0)
	if schema.Title != nil {
		content = appendPair(content, "title", nodeForString(*schema.Title))
	}
	if schema.ID != nil {
		content = appendPair(content, "id", nodeForString(*schema.ID))
	}
	if schema.Schema != nil {
		content = appendPair(content, "$schema", nodeForString(*schema.Schema))
	}
	if schema.Type != nil {
		content = appendPair(content, "type", schema.Type.nodeValue())
	}
	if schema.Items != nil {
		content = appendPair(content, "items", schema.Items.nodeValue())
	}
	if schema.Description != nil {
		content = appendPair(content, "description", nodeForString(*schema.Description))
	}
	if schema.Required != nil {
		content = appendPair(content, "required", nodeForStringArray(*schema.Required))
	}
	if schema.AdditionalProperties != nil {
		content = appendPair(content, "additionalProperties", schema.AdditionalProperties.nodeValue())
	}
	if schema.PatternProperties != nil {
		content = appendPair(content, "patternProperties", nodeForNamedSchemaArray(schema.PatternProperties))
	}
	if schema.Properties != nil {
		content = appendPair(content, "properties", nodeForNamedSchemaArray(schema.Properties))
	}
	if schema.Dependencies != nil {
		content = appendPair(content, "dependencies", nodeForNamedSchemaOrStringArray(schema.Dependencies))
	}
	if schema.Ref != nil {
		content = appendPair(content, "$ref", nodeForString(*schema.Ref))
	}
	if schema.MultipleOf != nil {
		content = appendPair(content, "multipleOf", schema.MultipleOf.nodeValue())
	}
	if schema.Maximum != nil {
		content = appendPair(content, "maximum", schema.Maximum.nodeValue())
	}
	if schema.ExclusiveMaximum != nil {
		content = appendPair(content, "exclusiveMaximum", nodeForBoolean(*schema.ExclusiveMaximum))
	}
	if schema.Minimum != nil {
		content = appendPair(content, "minimum", schema.Minimum.nodeValue())
	}
	if schema.ExclusiveMinimum != nil {
		content = appendPair(content, "exclusiveMinimum", nodeForBoolean(*schema.ExclusiveMinimum))
	}
	if schema.MaxLength != nil {
		content = appendPair(content, "maxLength", nodeForInt64(*schema.MaxLength))
	}
	if schema.MinLength != nil {
		content = appendPair(content, "minLength", nodeForInt64(*schema.MinLength))
	}
	if schema.Pattern != nil {
		content = appendPair(content, "pattern", nodeForString(*schema.Pattern))
	}
	if schema.AdditionalItems != nil {
		content = appendPair(content, "additionalItems", schema.AdditionalItems.nodeValue())
	}
	if schema.MaxItems != nil {
		content = appendPair(content, "maxItems", nodeForInt64(*schema.MaxItems))
	}
	if schema.MinItems != nil {
		content = appendPair(content, "minItems", nodeForInt64(*schema.MinItems))
	}
	if schema.UniqueItems != nil {
		content = appendPair(content, "uniqueItems", nodeForBoolean(*schema.UniqueItems))
	}
	if schema.MaxProperties != nil {
		content = appendPair(content, "maxProperties", nodeForInt64(*schema.MaxProperties))
	}
	if schema.MinProperties != nil {
		content = appendPair(content, "minProperties", nodeForInt64(*schema.MinProperties))
	}
	if schema.Enumeration != nil {
		content = appendPair(content, "enum", nodeForSchemaEnumArray(schema.Enumeration))
	}
	if schema.AllOf != nil {
		content = appendPair(content, "allOf", nodeForSchemaArray(*schema.AllOf))
	}
	if schema.AnyOf != nil {
		content = appendPair(content, "anyOf", nodeForSchemaArray(*schema.AnyOf))
	}
	if schema.OneOf != nil {
		content = appendPair(content, "oneOf", nodeForSchemaArray(*schema.OneOf))
	}
	if schema.Not != nil {
		content = appendPair(content, "not", schema.Not.nodeValue())
	}
	if schema.Definitions != nil {
		content = appendPair(content, "definitions", nodeForNamedSchemaArray(schema.Definitions))
	}
	if schema.Default != nil {
		// m = append(m, yaml.MapItem{Key: "default", Value: *schema.Default})
	}
	if schema.Format != nil {
		content = appendPair(content, "format", nodeForString(*schema.Format))
	}
	n.Content = content
	return n
=======
func namedSchemaArrayValue(array *[]*NamedSchema) interface{} {
	m2 := yaml.MapSlice{}
	for _, pair := range *(array) {
		var item2 yaml.MapItem
		item2.Key = pair.Name
		item2.Value = pair.Value.jsonValue()
		m2 = append(m2, item2)
	}
	return m2
}

func namedSchemaOrStringArrayValue(array *[]*NamedSchemaOrStringArray) interface{} {
	m2 := yaml.MapSlice{}
	for _, pair := range *(array) {
		var item2 yaml.MapItem
		item2.Key = pair.Name
		item2.Value = pair.Value.jsonValue()
		m2 = append(m2, item2)
	}
	return m2
}

func schemaEnumArrayValue(array *[]SchemaEnumValue) []interface{} {
	a := make([]interface{}, 0)
	for _, item := range *array {
		a = append(a, item.jsonValue())
	}
	return a
}

func schemaArrayValue(array *[]*Schema) []interface{} {
	a := make([]interface{}, 0)
	for _, item := range *array {
		a = append(a, item.jsonValue())
	}
	return a
}

func (schema *Schema) jsonValue() yaml.MapSlice {
	m := yaml.MapSlice{}
	if schema.Title != nil {
		m = append(m, yaml.MapItem{Key: "title", Value: *schema.Title})
	}
	if schema.ID != nil {
		m = append(m, yaml.MapItem{Key: "id", Value: *schema.ID})
	}
	if schema.Schema != nil {
		m = append(m, yaml.MapItem{Key: "$schema", Value: *schema.Schema})
	}
	if schema.Type != nil {
		m = append(m, yaml.MapItem{Key: "type", Value: schema.Type.jsonValue()})
	}
	if schema.Items != nil {
		m = append(m, yaml.MapItem{Key: "items", Value: schema.Items.jsonValue()})
	}
	if schema.Description != nil {
		m = append(m, yaml.MapItem{Key: "description", Value: *schema.Description})
	}
	if schema.Required != nil {
		m = append(m, yaml.MapItem{Key: "required", Value: *schema.Required})
	}
	if schema.AdditionalProperties != nil {
		m = append(m, yaml.MapItem{Key: "additionalProperties", Value: schema.AdditionalProperties.jsonValue()})
	}
	if schema.PatternProperties != nil {
		m = append(m, yaml.MapItem{Key: "patternProperties", Value: namedSchemaArrayValue(schema.PatternProperties)})
	}
	if schema.Properties != nil {
		m = append(m, yaml.MapItem{Key: "properties", Value: namedSchemaArrayValue(schema.Properties)})
	}
	if schema.Dependencies != nil {
		m = append(m, yaml.MapItem{Key: "dependencies", Value: namedSchemaOrStringArrayValue(schema.Dependencies)})
	}
	if schema.Ref != nil {
		m = append(m, yaml.MapItem{Key: "$ref", Value: *schema.Ref})
	}
	if schema.MultipleOf != nil {
		m = append(m, yaml.MapItem{Key: "multipleOf", Value: schema.MultipleOf.jsonValue()})
	}
	if schema.Maximum != nil {
		m = append(m, yaml.MapItem{Key: "maximum", Value: schema.Maximum.jsonValue()})
	}
	if schema.ExclusiveMaximum != nil {
		m = append(m, yaml.MapItem{Key: "exclusiveMaximum", Value: schema.ExclusiveMaximum})
	}
	if schema.Minimum != nil {
		m = append(m, yaml.MapItem{Key: "minimum", Value: schema.Minimum.jsonValue()})
	}
	if schema.ExclusiveMinimum != nil {
		m = append(m, yaml.MapItem{Key: "exclusiveMinimum", Value: schema.ExclusiveMinimum})
	}
	if schema.MaxLength != nil {
		m = append(m, yaml.MapItem{Key: "maxLength", Value: *schema.MaxLength})
	}
	if schema.MinLength != nil {
		m = append(m, yaml.MapItem{Key: "minLength", Value: *schema.MinLength})
	}
	if schema.Pattern != nil {
		m = append(m, yaml.MapItem{Key: "pattern", Value: *schema.Pattern})
	}
	if schema.AdditionalItems != nil {
		m = append(m, yaml.MapItem{Key: "additionalItems", Value: schema.AdditionalItems.jsonValue()})
	}
	if schema.MaxItems != nil {
		m = append(m, yaml.MapItem{Key: "maxItems", Value: *schema.MaxItems})
	}
	if schema.MinItems != nil {
		m = append(m, yaml.MapItem{Key: "minItems", Value: *schema.MinItems})
	}
	if schema.UniqueItems != nil {
		m = append(m, yaml.MapItem{Key: "uniqueItems", Value: *schema.UniqueItems})
	}
	if schema.MaxProperties != nil {
		m = append(m, yaml.MapItem{Key: "maxProperties", Value: *schema.MaxProperties})
	}
	if schema.MinProperties != nil {
		m = append(m, yaml.MapItem{Key: "minProperties", Value: *schema.MinProperties})
	}
	if schema.Enumeration != nil {
		m = append(m, yaml.MapItem{Key: "enum", Value: schemaEnumArrayValue(schema.Enumeration)})
	}
	if schema.AllOf != nil {
		m = append(m, yaml.MapItem{Key: "allOf", Value: schemaArrayValue(schema.AllOf)})
	}
	if schema.AnyOf != nil {
		m = append(m, yaml.MapItem{Key: "anyOf", Value: schemaArrayValue(schema.AnyOf)})
	}
	if schema.OneOf != nil {
		m = append(m, yaml.MapItem{Key: "oneOf", Value: schemaArrayValue(schema.OneOf)})
	}
	if schema.Not != nil {
		m = append(m, yaml.MapItem{Key: "not", Value: schema.Not.jsonValue()})
	}
	if schema.Definitions != nil {
		m = append(m, yaml.MapItem{Key: "definitions", Value: namedSchemaArrayValue(schema.Definitions)})
	}
	if schema.Default != nil {
		m = append(m, yaml.MapItem{Key: "default", Value: *schema.Default})
	}
	if schema.Format != nil {
		m = append(m, yaml.MapItem{Key: "format", Value: *schema.Format})
	}
	return m
>>>>>>> 79bfea2d (update vendor)
}

// JSONString returns a json representation of a schema.
func (schema *Schema) JSONString() string {
<<<<<<< HEAD
	node := schema.nodeValue()
	return Render(node)
=======
	info := schema.jsonValue()
	return Render(info)
>>>>>>> 79bfea2d (update vendor)
}
