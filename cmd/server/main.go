package main

import (
	"fmt"
	"os"

	v1 "github.com/nedson202/auth-manager/pkg/api/v1"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	if err := v1.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
