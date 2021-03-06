package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/radoondas/jmxproxybeat/beater"
)

func main() {
	err := beat.Run("jmxproxybeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
