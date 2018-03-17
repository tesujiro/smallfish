package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/tesujiro/smallfish/consumer/geo"
)

func main() {
	//csgeo.Run()
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s", err)
			os.Exit(1)
		}
	}()
	os.Exit(_main())
}

func _main() int {
	if envvar := os.Getenv("GOMAXPROCS"); envvar == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	csgeo.Run(ctx)

	return 0
}
