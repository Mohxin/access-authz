package main

import (
	"fmt"
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/volvo-cars/connect-access-control/internal/config"
	"github.com/volvo-cars/connect-access-control/internal/pkg/validator"
)

const (
	schemaDir string = "config/schema"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	v := validator.NewSchemaValidator(cfg.IAM.RootDir, schemaDir)
	errCode := 0

	results, err := v.Validate()
	if err != nil {
		println(err.Error())
		errCode = 1
	}

	for _, result := range results {
		if result.Valid() {
			fmt.Printf("[√] file://%s\n", result.FilePath)
			continue
		}

		fmt.Printf("[x] file://%s\n", result.FilePath)
		for i, err := range result.Errors() {
			fmt.Printf("	[%d] Message  	  :: %s\n", i, err.Message)
			fmt.Printf("	[%d] Field	  :: %s\n", i, err.Field)
			fmt.Printf("	[%d] Value	  :: %s\n", i, err.Value)
			fmt.Println("	------------------------------------------------")
			errCode = 1
		}
	}

	os.Exit(errCode)
}
