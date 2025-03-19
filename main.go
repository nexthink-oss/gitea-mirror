package main

import (
	"context"
	"os"

	"github.com/nexthink-oss/gitea-mirror/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.New().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
