package routes

import (
	"github.com/unkeyed/unkey/go/internal/services/auditlogs"
	"github.com/unkeyed/unkey/go/internal/services/caches"
	"github.com/unkeyed/unkey/go/internal/services/keys"
	"github.com/unkeyed/unkey/go/internal/services/ratelimit"
	"github.com/unkeyed/unkey/go/pkg/clickhouse"
	"github.com/unkeyed/unkey/go/pkg/clickhouse/schema"
	"github.com/unkeyed/unkey/go/pkg/db"
	"github.com/unkeyed/unkey/go/pkg/otel/logging"
	"github.com/unkeyed/unkey/go/pkg/vault"
	"github.com/unkeyed/unkey/go/pkg/zen/validation"
)

type EventBuffer interface {
	BufferApiRequest(schema.ApiRequestV1)
}

type Services struct {
	Logger         logging.Logger
	Database       db.Database
	Keys           keys.KeyService
	ClickHouse     clickhouse.ClickHouse
	Validator      *validation.Validator
	Ratelimit      ratelimit.Service
	Auditlogs      auditlogs.AuditLogService
	Caches         caches.Caches
	Vault          *vault.Service
	ChproxyEnabled bool
	ChproxyToken   string
}
