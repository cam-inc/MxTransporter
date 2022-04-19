package file

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func Test_FileExporter(t *testing.T) {
	t.Run("File Exporter running", func(t *testing.T) {
		New(&ExporterConfig{})
		e := New(&ExporterConfig{
			LogType:         "debug",
			ChangeStreamKey: "changeStreamKey",
		})

		e.Export(context.Background(), primitive.M{
			"_id": "xxxxxxxxxxxx",
		})
	})
}