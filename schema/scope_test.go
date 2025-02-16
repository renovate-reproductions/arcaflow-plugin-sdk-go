package schema_test

import (
	"encoding/json"
	"go.arcalot.io/assert"
	"testing"

	"go.flow.arcalot.io/pluginsdk/schema"
)

type scopeTestObjectEmpty struct {
}

type scopeTestObjectB struct {
	C string `json:"c"`
}

type scopeTestObjectA struct {
	B scopeTestObjectB `json:"b"`
}

type scopeTestObjectAPtr struct {
	B *scopeTestObjectB `json:"b"`
}

var scopeTestObjectEmptySchema = schema.NewScopeSchema(
	schema.NewStructMappedObjectSchema[scopeTestObjectEmpty](
		"scopeTestObjectEmpty",
		map[string]*schema.PropertySchema{},
	),
)
var scopeTestObjectEmptySchemaRenamed = schema.NewScopeSchema(
	schema.NewStructMappedObjectSchema[scopeTestObjectEmpty](
		"scopeTestObjectEmptyRenamed",
		map[string]*schema.PropertySchema{},
	),
)

var scopeTestObjectCStrSchema = schema.NewScopeSchema(
	schema.NewObjectSchema(
		"scopeTestObjectC",
		map[string]*schema.PropertySchema{
			"d": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
)
var scopeTestObjectCIntSchema = schema.NewScopeSchema(
	schema.NewObjectSchema(
		"scopeTestObjectC",
		map[string]*schema.PropertySchema{
			"d": schema.NewPropertySchema(
				schema.NewIntSchema(nil, nil, nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
)

var scopeTestObjectASchema = schema.NewScopeSchema(
	schema.NewObjectSchema(
		"scopeTestObjectA",
		map[string]*schema.PropertySchema{
			"b": schema.NewPropertySchema(
				schema.NewRefSchema("scopeTestObjectB", nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
	schema.NewObjectSchema(
		"scopeTestObjectB",
		map[string]*schema.PropertySchema{
			"c": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
)

var scopeTestObjectAType = schema.NewScopeSchema(
	schema.NewStructMappedObjectSchema[scopeTestObjectA](
		"scopeTestObjectA",
		map[string]*schema.PropertySchema{
			"b": schema.NewPropertySchema(
				schema.NewRefSchema("scopeTestObjectB", nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
	schema.NewStructMappedObjectSchema[scopeTestObjectB](
		"scopeTestObjectB",
		map[string]*schema.PropertySchema{
			"c": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
)

var scopeTestObjectATypePtr = schema.NewScopeSchema(
	schema.NewStructMappedObjectSchema[*scopeTestObjectAPtr](
		"scopeTestObjectA",
		map[string]*schema.PropertySchema{
			"b": schema.NewPropertySchema(
				schema.NewRefSchema("scopeTestObjectB", nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
	schema.NewStructMappedObjectSchema[*scopeTestObjectB](
		"scopeTestObjectB",
		map[string]*schema.PropertySchema{
			"c": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, nil),
				nil,
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
)

func TestScopeConstructor(t *testing.T) {
	assert.Equals(t, scopeTestObjectASchema.TypeID(), schema.TypeIDScope)
	assert.Equals(t, scopeTestObjectAType.TypeID(), schema.TypeIDScope)
}

func TestUnserialization(t *testing.T) {
	// Test unserialization of composition of two objects
	data := `{"b":{"c": "Hello world!"}}`
	var input any
	assert.NoError(t, json.Unmarshal([]byte(data), &input))

	result, err := scopeTestObjectAType.Unserialize(input)
	assert.NoError(t, err)
	assert.InstanceOf[scopeTestObjectA](t, result.(scopeTestObjectA))
	assert.Equals(t, result.(scopeTestObjectA).B.C, "Hello world!")

	// Now as a ptr
	resultPtr, err := scopeTestObjectATypePtr.Unserialize(input)
	assert.NoError(t, err)
	assert.InstanceOf[*scopeTestObjectAPtr](t, resultPtr.(*scopeTestObjectAPtr))
	assert.Equals(t, resultPtr.(*scopeTestObjectAPtr).B.C, "Hello world!")

	// Test empty object
	data = `{}`
	assert.NoError(t, json.Unmarshal([]byte(data), &input))
	result, err = scopeTestObjectEmptySchema.Unserialize(input)
	assert.NoError(t, err)
	assert.InstanceOf[scopeTestObjectEmpty](t, result.(scopeTestObjectEmpty))
}

func TestValidation(t *testing.T) {
	// Note: The scopeTestObject var used must be NewStructMappedObjectSchema,
	// or else it will be a dict instead of a struct, causing problems.
	// Test composition of two objects
	err := scopeTestObjectAType.Validate(scopeTestObjectA{
		scopeTestObjectB{
			"Hello world!",
		},
	})
	assert.NoError(t, err)

	// Test empty scope object
	err = scopeTestObjectEmptySchema.Validate(scopeTestObjectEmpty{})
	assert.NoError(t, err)
}
func TestCompatibilityValidationWithData(t *testing.T) {
	err := scopeTestObjectAType.ValidateCompatibility(map[string]any{
		"b": map[string]any{
			"c": "Hello world!",
		},
	})
	assert.NoError(t, err)

	// Replace the actual value with a schema
	err = scopeTestObjectAType.ValidateCompatibility(map[string]any{
		"b": map[string]any{
			"c": schema.NewStringSchema(nil, nil, nil),
		},
	})
	assert.NoError(t, err)

	// Test empty scope object
	// The ValidateCompatibility method should behave like Validate when data is passed in
	err = scopeTestObjectEmptySchema.ValidateCompatibility(map[string]any{})
	assert.NoError(t, err)
}

func TestCompatibilityValidationWithSchema(t *testing.T) {
	// Note: The scopeTestObject var used must be NewStructMappedObjectSchema,
	// or else it will be a dict instead of a struct, causing problems.
	// Test composition of two objects

	// Note: Doesn't support the non-pointer, dereferenced version of the scope type.
	err := scopeTestObjectAType.ValidateCompatibility(scopeTestObjectAType)
	assert.NoError(t, err)

	// Test empty scope object
	// Note: Doesn't support the non-pointer version.
	err = scopeTestObjectEmptySchema.ValidateCompatibility(scopeTestObjectEmptySchema)
	assert.NoError(t, err)

	// Now mismatched
	err = scopeTestObjectAType.ValidateCompatibility(scopeTestObjectEmptySchema)
	assert.Error(t, err)
	err = scopeTestObjectEmptySchema.ValidateCompatibility(scopeTestObjectAType)
	assert.Error(t, err)

	// Similar, but with a simple difference
	// Mismatching IDs
	err = scopeTestObjectEmptySchema.ValidateCompatibility(scopeTestObjectEmptySchemaRenamed)
	assert.Error(t, err)
	err = scopeTestObjectEmptySchemaRenamed.ValidateCompatibility(scopeTestObjectEmptySchema)
	assert.Error(t, err)
	// Mismatching type in one field, but with the field ID matching
	err = scopeTestObjectCStrSchema.ValidateCompatibility(scopeTestObjectCIntSchema)
	assert.Error(t, err)
	err = scopeTestObjectCIntSchema.ValidateCompatibility(scopeTestObjectCStrSchema)
	assert.Error(t, err)
}

func TestSerialization(t *testing.T) {
	serialized, err := scopeTestObjectAType.Serialize(scopeTestObjectA{
		scopeTestObjectB{
			"Hello world!",
		},
	})
	assert.NoError(t, err)
	assert.Equals(t, serialized.(map[string]any)["b"].(map[string]any)["c"].(string), "Hello world!")
}

func TestSelfSerialization(t *testing.T) {
	serializedScope, err := scopeTestObjectAType.SelfSerialize()
	assert.NoError(t, err)
	serializedScopeMap := serializedScope.(map[string]any)
	if serializedScopeMap["root"] != "scopeTestObjectA" {
		t.Fatalf("Unexpected root object: %s", serializedScopeMap["root"])
	}
}
