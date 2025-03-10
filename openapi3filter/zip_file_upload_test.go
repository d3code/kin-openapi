package openapi3filter_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/d3code/kin-openapi/openapi3"
	"github.com/d3code/kin-openapi/openapi3filter"
	"github.com/d3code/kin-openapi/routers/gorillamux"
)

func TestValidateZipFileUpload(t *testing.T) {
	const spec = `
openapi: 3.0.0
info:
  title: 'Validator'
  version: 0.0.1
paths:
  /test:
    post:
      requestBody:
        required: true
        content:
          multipart/form-data:
            schema:
              type: object
              required:
                - file
              properties:
                file:
                  type: string
                  format: binary
      responses:
        '200':
          description: Created
`

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData([]byte(spec))
	require.NoError(t, err)

	err = doc.Validate(loader.Context)
	require.NoError(t, err)

	router, err := gorillamux.NewRouter(doc)
	require.NoError(t, err)

	tests := []struct {
		zipData []byte
		wantErr bool
	}{
		{
			[]byte{
				0x50, 0x4b, 0x03, 0x04, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x7d, 0x23, 0x56, 0xcd, 0xfd, 0x67, 0xf8, 0x07, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x09, 0x00, 0x1c, 0x00, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x74, 0x78, 0x74, 0x55, 0x54, 0x09, 0x00, 0x03, 0xac, 0xce, 0xb3, 0x63, 0xaf, 0xce, 0xb3, 0x63, 0x75, 0x78, 0x0b, 0x00, 0x01, 0x04, 0xf7, 0x01, 0x00, 0x00, 0x04, 0x14, 0x00, 0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x0a, 0x50, 0x4b, 0x01, 0x02, 0x1e, 0x03, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x7d, 0x23, 0x56, 0xcd, 0xfd, 0x67, 0xf8, 0x07, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x09, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0xa4, 0x81, 0x00, 0x00, 0x00, 0x00, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x74, 0x78, 0x74, 0x55, 0x54, 0x05, 0x00, 0x03, 0xac, 0xce, 0xb3, 0x63, 0x75, 0x78, 0x0b, 0x00, 0x01, 0x04, 0xf7, 0x01, 0x00, 0x00, 0x04, 0x14, 0x00, 0x00, 0x00, 0x50, 0x4b, 0x05, 0x06, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x4f, 0x00, 0x00, 0x00, 0x4a, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			false,
		},
		{
			[]byte{
				0x50, 0x4b, 0x05, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}, // No entry
			true,
		},
	}
	for _, tt := range tests {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		{ // Add file data
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", `form-data; name="file"; filename="hello.zip"`)
			h.Set("Content-Type", "application/zip")

			fw, err := writer.CreatePart(h)
			require.NoError(t, err)
			_, err = io.Copy(fw, bytes.NewReader(tt.zipData))

			require.NoError(t, err)
		}

		writer.Close()

		req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(body.Bytes()))
		require.NoError(t, err)

		req.Header.Set("Content-Type", writer.FormDataContentType())

		route, pathParams, err := router.FindRoute(req)
		require.NoError(t, err)

		if err = openapi3filter.ValidateRequestBody(
			context.Background(),
			&openapi3filter.RequestValidationInput{
				Request:    req,
				PathParams: pathParams,
				Route:      route,
			},
			route.Operation.RequestBody.Value,
		); err != nil {
			if !tt.wantErr {
				t.Errorf("got %v", err)
			}
			continue
		}
		if tt.wantErr {
			t.Errorf("want err")
		}
	}
}
