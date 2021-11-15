package resume_token

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mxtransporter/config"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type MockGeneralConfig struct {
	config.GeneralConfigIf
	FakePersistentVolume func() (string, error)
}

func (m *MockGeneralConfig) FetchPersistentVolumeDir() (string, error) {
	return m.FakePersistentVolume()
}

func Test_SaveResumeToken(t *testing.T) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal("Failed to load location time.")
	}

	nowTime := time.Now().In(jst)
	file := nowTime.Format("2006/01/02/2006-01-02.dat")

	rtMap := primitive.M{"_data": "00000"}

	cases := []struct {
		rt       primitive.M
		function config.GeneralConfigIf
	}{
		{
			rt: rtMap,
			function: &MockGeneralConfig{
				FakePersistentVolume: func() (string, error) {
					return "", nil
				},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := NewGeneralConfig(tt.function).SaveResumeToken(tt.rt); err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			rtByte, err := os.ReadFile(file)
			if err != nil {
				t.Fatal("Failed to read file saved test resume token in.")
			}

			if e, a := rtMap["_data"], string(rtByte); !reflect.DeepEqual(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestMain(m *testing.M) {
	status := m.Run()

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Failed to load location time.")
	}

	nowTime := time.Now().In(jst)

	err = os.RemoveAll(nowTime.Format("2006"))
	if err != nil {
		fmt.Println("The unnecessary file could not be deleted.")
	}

	os.Exit(status)
}
