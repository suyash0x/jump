package jump

import (
	"log"
	"os"
)

func FatalOutError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
