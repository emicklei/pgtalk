package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type TableFilter struct {
	includes, excludes []*regexp.Regexp
}

func NewTableFilter(includes, excludes string) TableFilter {
	f := TableFilter{}
	inc := strings.Split(includes, ",")
	exc := strings.Split(excludes, ",")
	for _, each := range inc {
		entry := strings.TrimSpace(each)
		if !strings.Contains(entry, "*") {
			entry = fmt.Sprintf("^%s$", entry)
		}
		if len(entry) > 0 {
			f.includes = append(f.includes, regexp.MustCompile(entry))
		}
	}
	for _, each := range exc {
		entry := strings.TrimSpace(each)
		if !strings.Contains(entry, "*") {
			entry = fmt.Sprintf("^%s$", entry)
		}
		if len(entry) > 0 {
			f.excludes = append(f.excludes, regexp.MustCompile(entry))
		}
	}
	return f
}

func (f TableFilter) Includes(name string) bool {
	if len(f.includes) > 0 {
		// must be in includes
		included := false
		for _, each := range f.includes {
			if each.MatchString(name) {
				included = true
				continue
			}
		}
		if !included {
			if *oVerbose {
				log.Println("[skip] filters does not include:", name)
			}
			return false
		}
	}
	// must not be in excludes
	for _, each := range f.excludes {
		if each.MatchString(name) {
			if *oVerbose {
				log.Println("[skip] filters does exclude:", name)
			}
			return false
		}
	}
	return true
}
