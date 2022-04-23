package pluginimports

// We need to empty import all enabled plugins so that they will be linked into
// any hook binary.
import (
	_ "github.com/cheld/miniprow/pkg/hook/plugins/triggers/github" // Import all enabled plugins.
	_ "github.com/cheld/miniprow/pkg/hook/plugins/triggers/http"
)
