// Prints the DDL for the GORM models; atlas.hcl uses it as the desired schema.
package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/eunice99x/pixely/kyc/internal/model"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&model.Image{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
