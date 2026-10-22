package main

import (
	"fmt"
	"strings"
)

func getParentDomain(domain string) string {
	parts := strings.Split(domain, ".")
	parents := parts[1:]
	parentDomain := strings.Join(parents, ".")
	return parentDomain
}

func main() {
	url := "www.google.com."
	parent := getParentDomain(url)
	fmt.Println(url, parent)
}
