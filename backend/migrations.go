package main

import "embed"

//go:embed all:database/migrations
var migrationsFS embed.FS
