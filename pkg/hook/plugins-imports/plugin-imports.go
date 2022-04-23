package pluginimports

// We need to empty import all enabled plugins so that they will be linked into
// any hook binary.
import (

	// triggers
	_ "github.com/cheld/miniprow/pkg/hook/plugins/triggers/github" // Import all enabled plugins.
	_ "github.com/cheld/miniprow/pkg/hook/plugins/triggers/http"

	// actions
	_ "github.com/cheld/miniprow/pkg/hook/plugins/actions/debug"
	_ "github.com/cheld/miniprow/pkg/hook/plugins/actions/github"
	_ "github.com/cheld/miniprow/pkg/hook/plugins/actions/http"
	_ "github.com/cheld/miniprow/pkg/hook/plugins/actions/misc"
)
