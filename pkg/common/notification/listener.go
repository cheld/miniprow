package notification

import (
	"context"

	"github.com/cheld/miniprow/pkg/common/core"
)

type Listener func(*core.Event, core.Tenant, context.Context)
