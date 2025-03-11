package migrations

import "embed"

//go:embed *.sql
var AuthEmbedFS embed.FS
