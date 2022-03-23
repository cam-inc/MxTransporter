//go:build test
// +build test

package resume_token

/*
func Test_File_SaveResumeToken(t *testing.T) {
	var l *zap.SugaredLogger
	logConfig := config.LogConfig()
	l = logger.New(logConfig)

	rt := "00000"

	currentDir, _ := os.Getwd()
	ctx := context.Background()
	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Failed to fetch RESUME_TOKEN_VOLUME_DIR value.",
			runner: func(t *testing.T) {
				_, err := New(ctx, l, "file")
				if err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Pass to save resume token in file.",
			runner: func(t *testing.T) {
				if err := setEnv(constant.RESUME_TOKEN_VOLUME_DIR, currentDir); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_DIR environment variables.")
				}

				if err := setEnv(constant.MONGODB_COLLECTION, "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}

				resumeTokenImpl, err := New(ctx, l, "file")
				if err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				env := resumeTokenImpl.Env()
				if env == "" {
					t.Fatalf("Failed to get environment variables for resume tokens settings.")
				}
				fmt.Printf("env %s\n", env)

				if err := resumeTokenImpl.SaveResumeToken(ctx, rt); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				resumeToken := resumeTokenImpl.ReadResumeToken(ctx)
				if resumeToken == "" {
					t.Fatal("Failed to read file saved test resume token in.")
				}
				fmt.Printf("resumeToken %s\n", resumeToken)
				envMap := map[string]string{}
				if err := json.Unmarshal([]byte(env), &envMap); err != nil {
					t.Fatal(err)
				}

				// remove resumeToken file for test
				os.RemoveAll(fmt.Sprintf("%s/%s", envMap["volume_path"], envMap["file_name"]))

				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_DIR); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_DIR environment variables.")
				}
				if err := unsetEnv(constant.MONGODB_COLLECTION); err != nil {
					t.Fatalf("Failed to unset file MONGODB_COLLECTION environment variables.")
				}

			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}

//func setEnv(key, value string) error {
//	return os.Setenv(key, value)
//}
//
//func unsetEnv(key string) error {
//	return os.Unsetenv(key)
//}
//
*/
