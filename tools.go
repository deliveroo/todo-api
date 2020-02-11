// +build tools

// Package tools records build-time dependencies that aren't used by the
// library itself, but are tracked by go mod and required to lint and
// build the project.
package tools

import (
	_ "github.com/cortesi/modd/cmd/modd"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/johngibb/migrate/cmd/migrate"
	_ "github.com/joho/godotenv/cmd/godotenv"
	_ "github.com/jstemmer/go-junit-report"
)
