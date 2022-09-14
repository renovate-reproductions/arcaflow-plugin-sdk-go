package schema_test

import (
	"testing"

	"go.flow.arcalot.io/pluginsdk/schema"
)

type serializationTestCase[T any] struct {
	SerializedValue         any
	ExpectError             bool
	ExpectUnserializedValue T
	ExpectedSerializedValue any
}

func performSerializationTest[T any](
	t *testing.T,
	typeUnderTest schema.AbstractType[T],
	testCases map[string]serializationTestCase[T],
	compareUnserialized func(a T, b T) bool,
	compareSerialized func(a any, b any) bool,
) {
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			unserialized, err := typeUnderTest.Unserialize(tc.SerializedValue)
			if err != nil {
				if tc.ExpectError {
					return
				}
				t.Fatal(err)
			}
			if err := typeUnderTest.Validate(unserialized); err != nil {
				t.Fatal(err)
			}
			if !compareUnserialized(unserialized, tc.ExpectUnserializedValue) {
				t.Fatalf(
					"Unexpected unserialized value, expected: %v, got: %v",
					tc.ExpectUnserializedValue,
					unserialized,
				)
			}
			serialized, err := typeUnderTest.Serialize(unserialized)
			if err != nil {
				t.Fatal(err)
			}
			if !compareSerialized(serialized, tc.ExpectedSerializedValue) {
				t.Fatalf(
					"Serialized value mismatch, expected: %v (%T), got: %v (%T)",
					tc.ExpectedSerializedValue,
					tc.ExpectedSerializedValue,
					serialized,
					serialized,
				)
			}
		})
	}
}

func assertEqual[T comparable](t *testing.T, got T, expected T) {
	t.Helper()
	if expected != got {
		t.Fatalf("Mismatch, expected: %v, got: %v", expected, got)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error, no error returned")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error returned: %v", err)
	}
}

func assertError2(t *testing.T) func(_ any, err error) {
	t.Helper()
	return func(_ any, err error) {
		t.Helper()
		if err == nil {
			t.Fatalf("Expected error, no error returned")
		}
	}
}

func assertNoError2(t *testing.T) func(_ any, err error) {
	t.Helper()
	return func(_ any, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("Unexpected error returned: %v", err)
		}
	}
}
