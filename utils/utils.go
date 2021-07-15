package utils

import "log"

// when err isn't nil
// then print err
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
