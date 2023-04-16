package main

import (
	"fmt"
	"testing"
)

func TestReportTypes(t *testing.T) {
	for k := range pgMappings {
		fmt.Println("-", k)
	}
}
