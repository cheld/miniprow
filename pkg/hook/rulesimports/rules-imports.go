package rulesimports

// We need to empty import all enabled plugins so that they will be linked into
// any hook binary.
import (

	// triggers
	_ "github.com/cheld/miniprow/pkg/hook/rules/filters/github" // Import all enabled plugins.
	_ "github.com/cheld/miniprow/pkg/hook/rules/filters/http"

	// actions
	_ "github.com/cheld/miniprow/pkg/hook/rules/handlers/github"
	_ "github.com/cheld/miniprow/pkg/hook/rules/handlers/http"
	_ "github.com/cheld/miniprow/pkg/hook/rules/handlers/misc"
)
