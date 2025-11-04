package handlers

import (
	"burrowfs/api/schemas/webdav"
	database "burrowfs/core/db"
	"burrowfs/core/db/models/files"
	"burrowfs/core/types"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// HandleLock processes a LOCK request for a WebDAV resource.
func HandleLock(user *types.InternalUser, path *types.Path, depthHeader string, timeoutHeader string, body webdav.LockInfo) (int, *webdav.LockResponse, string) {
	dbConn, err := database.Open()
	defer dbConn.Close()
	if err != nil {
		return 500, nil, ""
	}
	resourceToLock, err := files.GetFileByPath(dbConn, user.Id, path.Clean)
	if err != nil || resourceToLock == nil {
		return 404, nil, ""
	}

	if resourceToLock.IsLocked() {
		return 423, nil, ""
	}

	depth := normalizeDepth(depthHeader)
	timeoutDur, timeoutOut := parseTimeout(timeoutHeader, 30*time.Minute, 24*time.Hour)

	token := fmt.Sprintf("opaquelocktoken:%s", uuid.New().String())

	if err := files.LockFileWithToken(dbConn, *resourceToLock.FileID, user.Id, token, timeoutDur); err != nil {
		return 500, nil, ""
	}

	scope := webdav.ScopeFromLockInfo(body)
	owner := webdav.OwnerFromLockInfo(body)
	active := webdav.NewActiveLock(depth, timeoutOut, token, path.Clean, scope, owner)
	resp := webdav.NewLockResponse(active)

	lockTokenHeader := fmt.Sprintf("<%s>", token)

	return 200, &resp, lockTokenHeader
}

func normalizeDepth(v string) string {
	switch v {
	case "0":
		return "0"
	default:
		return "infinity"
	}
}

func parseTimeout(header string, def time.Duration, max time.Duration) (time.Duration, string) {
	h := strings.TrimSpace(strings.ToLower(header))
	if h == "" || h == "infinite" {
		return def, fmt.Sprintf("Second-%d", int(def/time.Second))
	}
	if strings.HasPrefix(h, "second-") {
		secStr := strings.TrimPrefix(h, "second-")
		if n, err := time.ParseDuration(secStr + "s"); err == nil {
			if n <= 0 {
				return def, fmt.Sprintf("Second-%d", int(def/time.Second))
			}
			if n > max {
				n = max
			}
			return n, fmt.Sprintf("Second-%d", int(n/time.Second))
		}
	}
	return def, fmt.Sprintf("Second-%d", int(def/time.Second))
}
