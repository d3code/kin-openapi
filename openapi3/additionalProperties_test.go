package openapi3_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/oasdiff/yaml3"
	"github.com/stretchr/testify/require"

	"github.com/d3code/kin-openapi/openapi3"
)

func TestMarshalAdditionalProperties(t *testing.T) {
	ctx := context.Background()
	data, err := os.ReadFile("testdata/test.openapi.additionalproperties.yml")
	require.NoError(t, err)

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	spec, err := loader.LoadFromData(data)
	require.NoError(t, err)

	err = spec.Validate(ctx)
	require.NoError(t, err)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err = enc.Encode(spec)
	require.NoError(t, err)

	// Load the doc from the serialized yaml.
	spec2, err := loader.LoadFromData(buf.Bytes())
	require.NoError(t, err)

	err = spec2.Validate(ctx)
	require.NoError(t, err)
}
