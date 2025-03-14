package openapi3_test

import (
	"testing"

	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"

	"github.com/d3code/kin-openapi/openapi3"
)

func TestIssue689(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name     string
		schema   *openapi3.Schema
		value    map[string]any
		opts     []openapi3.SchemaValidationOption
		checkErr require.ErrorAssertionFunc
	}{
		// read-only
		{
			name: "read-only property succeeds when read-only validation is disabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, ReadOnly: null.BoolFrom(true)}}),
			value: map[string]any{"foo": true},
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsRequest(),
				openapi3.DisableReadOnlyValidation()},
			checkErr: require.NoError,
		},
		{
			name: "non read-only property succeeds when read-only validation is disabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, ReadOnly: null.BoolFrom(false)}}),
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsRequest()},
			value:    map[string]any{"foo": true},
			checkErr: require.NoError,
		},
		{
			name: "read-only property fails when read-only validation is enabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, ReadOnly: null.BoolFrom(true)}}),
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsRequest()},
			value:    map[string]any{"foo": true},
			checkErr: require.Error,
		},
		{
			name: "non read-only property succeeds when read-only validation is enabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, ReadOnly: null.BoolFrom(false)}}),
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsRequest()},
			value:    map[string]any{"foo": true},
			checkErr: require.NoError,
		},
		// write-only
		{
			name: "write-only property succeeds when write-only validation is disabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, WriteOnly: null.BoolFrom(true)}}),
			value: map[string]any{"foo": true},
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsResponse(),
				openapi3.DisableWriteOnlyValidation()},
			checkErr: require.NoError,
		},
		{
			name: "non write-only property succeeds when write-only validation is disabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, WriteOnly: null.BoolFrom(false)}}),
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsResponse()},
			value:    map[string]any{"foo": true},
			checkErr: require.NoError,
		},
		{
			name: "write-only property fails when write-only validation is enabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, WriteOnly: null.BoolFrom(true)}}),
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsResponse()},
			value:    map[string]any{"foo": true},
			checkErr: require.Error,
		},
		{
			name: "non write-only property succeeds when write-only validation is enabled",
			schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"foo": {Type: &openapi3.Types{"boolean"}, WriteOnly: null.BoolFrom(false)}}),
			opts: []openapi3.SchemaValidationOption{
				openapi3.VisitAsResponse()},
			value:    map[string]any{"foo": true},
			checkErr: require.NoError,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.schema.VisitJSON(test.value, test.opts...)
			test.checkErr(t, err)
		})
	}
}
