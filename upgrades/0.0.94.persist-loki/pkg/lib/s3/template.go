package s3

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

func generateTemplate(bucketName string) (string, error) {
	buf := bytes.Buffer{}

	t, err := template.New("bucket").Parse(rawBucketTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(&buf, struct {
		BucketName string
	}{
		BucketName: bucketName,
	})
	if err != nil {
		return "", fmt.Errorf("interpolating template: %w", err)
	}

	return buf.String(), nil
}

//go:embed bucket-template.yaml
var rawBucketTemplate string
