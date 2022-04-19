package file

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
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

	t.Run("Marshal and Unmarshal", func(t *testing.T) {
		ts := timestamp{
			Time: time.Now(),
		}
		data, err := json.Marshal(ts)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println(string(data))
		}

		ts2 := timestamp{}
		err = json.Unmarshal([]byte(`{"T":1650366482, "I":2}`), &ts2)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%v\n", ts2)
	})
}
