package master

import (
	"log"
	"regexp"
)

type Fingerprint string

func (f *Fingerprint) Valid() bool {
	matched, err := regexp.MatchString("^[a-f0-9]{40}$", string(*f))
	if err != nil {
		log.Fatal("Error while matching regex: %s", err)
	}
	return matched
}
