package docs

import (
	"embed"
)

// FS embeds all documentation files including the interactive Swagger UI and YAML specifications.
//
//go:embed index.html swagger.yaml modules/*.yaml
var FS embed.FS
