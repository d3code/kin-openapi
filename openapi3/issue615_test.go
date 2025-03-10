package openapi3_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/d3code/kin-openapi/openapi3"
)

func TestIssue615(t *testing.T) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err := loader.LoadFromFile("testdata/recursiveRef/issue615.yml")
	require.NoError(t, err)

	err = doc.Validate(loader.Context)
	require.NoError(t, err)
}
