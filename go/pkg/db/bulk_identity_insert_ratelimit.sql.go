// Code generated by sqlc bulk insert plugin. DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"strings"
)

// bulkInsertIdentityRatelimit is the base query for bulk insert
const bulkInsertIdentityRatelimit = `INSERT INTO ` + "`" + `ratelimits` + "`" + ` ( id, workspace_id, identity_id, name, ` + "`" + `limit` + "`" + `, duration, created_at, auto_apply ) VALUES %s ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    ` + "`" + `limit` + "`" + ` = VALUES(` + "`" + `limit` + "`" + `),
    duration = VALUES(duration),
    auto_apply = VALUES(auto_apply),
    updated_at = VALUES(created_at)`

// InsertIdentityRatelimits performs bulk insert in a single query
func (q *BulkQueries) InsertIdentityRatelimits(ctx context.Context, db DBTX, args []InsertIdentityRatelimitParams) error {

	if len(args) == 0 {
		return nil
	}

	// Build the bulk insert query
	valueClauses := make([]string, len(args))
	for i := range args {
		valueClauses[i] = "( ?, ?, ?, ?, ?, ?, ?, ? )"
	}

	bulkQuery := fmt.Sprintf(bulkInsertIdentityRatelimit, strings.Join(valueClauses, ", "))

	// Collect all arguments
	var allArgs []any
	for _, arg := range args {
		allArgs = append(allArgs, arg.ID)
		allArgs = append(allArgs, arg.WorkspaceID)
		allArgs = append(allArgs, arg.IdentityID)
		allArgs = append(allArgs, arg.Name)
		allArgs = append(allArgs, arg.Limit)
		allArgs = append(allArgs, arg.Duration)
		allArgs = append(allArgs, arg.CreatedAt)
		allArgs = append(allArgs, arg.AutoApply)
	}

	// Execute the bulk insert
	_, err := db.ExecContext(ctx, bulkQuery, allArgs...)
	return err
}
