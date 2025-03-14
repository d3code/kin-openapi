package openapi3_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/oasdiff/yaml3"
	"github.com/stretchr/testify/require"

	"github.com/d3code/kin-openapi/openapi3"
)

func TestIssue241(t *testing.T) {
	data, err := os.ReadFile("testdata/issue241.yml")
	require.NoError(t, err)

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	spec, err := loader.LoadFromData(data)
	require.NoError(t, err)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err = enc.Encode(spec)
	require.NoError(t, err)
	require.Equal(t, string(data), buf.String())
}
