package goclamav_test

import (
	"fmt"
	"log"
	"testing"

	clamav "github.com/ca110us/go-clamav"
)

func TestInit(t *testing.T) {
	clamavInstance := new(clamav.Clamav)
	err := clamavInstance.Init(clamav.SCAN_OPTIONS{})

	if err != nil {
		t.Errorf("Error with init: %v", err)
	}

	signo, err := clamavInstance.LoadDB("/var/lib/clamav", uint(clamav.CL_DB_DIRECTORY))
	if err != nil {
		t.Errorf("Error loading default db: %v", err)
	}
	log.Println("db load succeed:", signo)
}

func ExampleRetver() {
	fmt.Println(clamav.Retver())
}

func ExampleRetdbdir() {
	fmt.Println(clamav.Retdbdir())

	// Output:
	// /var/lib/clamav
}
