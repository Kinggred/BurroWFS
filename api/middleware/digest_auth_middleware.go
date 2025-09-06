package middleware

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/utils"
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	realmName     = "BurrowFS"
	nonceValidity = 5 * time.Minute
	opaqueValue   = "BurrowFS-WebDAV-Compatible"
)

var (
	nonceStore sync.Map
)

func GenerateNonce() string {
	hash := md5.Sum([]byte(time.Now().String()))
	return hex.EncodeToString(hash[:])
}

func hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func IsWebDAVMethod(method string) bool {
	webDAVMethods := map[string]bool{
		"PROPFIND":  true,
		"PROPPATCH": true,
		"MKCOL":     true,
		"COPY":      true,
		"MOVE":      true,
		"LOCK":      true,
		"UNLOCK":    true,
	}
	return webDAVMethods[method]
}

func ValidateDigest(username, uri, method, givenResponse, givenNonce string, dbConn *db.DB) bool {
	user, err := models.GetUserByEmail(dbConn, username)
	if err != nil || user == nil {
		return false
	}

	password, err := utils.DecodePassword(user.Password)
	if err != nil {
		return false
	}

	ha1 := hash([]byte(username + ":" + realmName + ":" + password))
	ha2 := hash([]byte(method + ":" + uri))

	// Simple digest authentication (RFC 2069)
	expectedResponse := hash([]byte(ha1 + ":" + givenNonce + ":" + ha2))

	return strings.ToLower(expectedResponse) == strings.ToLower(givenResponse)
}

func ParseDigestHeader(authHeader string) map[string]string {
	digestFields := make(map[string]string)

	for _, field := range strings.Split(authHeader, ",") {
		parts := strings.SplitN(strings.TrimSpace(field), "=", 2)
		if len(parts) == 2 {
			key := strings.Trim(parts[0], `" `)
			value := strings.Trim(parts[1], `" `)
			digestFields[key] = value
		}
	}
	return digestFields
}

func extractUserFromAuthHeader(digestFields map[string]string) schemas.UserResponse {
	username := digestFields["username"]

	doConn, err := db.Open()
	defer doConn.Close()
	if err != nil {
		return schemas.UserResponse{}
	}

	user, err := models.GetUserByEmail(doConn, username)
	if err != nil || user == nil {
		return schemas.UserResponse{}
	}

	return schemas.UserResponse{
		UserId:    user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.ModifiedAt.Format(time.RFC3339),
	}
}

func DigestAuthMiddleware(next http.Handler) http.Handler {
	/*
		WebDAV clients require RFC 2617 support with MD5 algorithm.
		This middleware supports both standard HTTP methods and WebDAV-specific methods.
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Special handling for OPTIONS method which is often used by WebDAV clients for discovery
		if r.Method == "OPTIONS" {
			// Add WebDAV headers
			w.Header().Set("DAV", "1, 2")
			w.Header().Set("MS-Author-Via", "DAV")
			w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST, PUT, DELETE, PROPFIND, PROPPATCH, MKCOL, COPY, MOVE, LOCK, UNLOCK")

			// If this is a preflight request, handle it and return
			if r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, POST, PUT, DELETE, PROPFIND, PROPPATCH, MKCOL, COPY, MOVE, LOCK, UNLOCK")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Depth, Destination, Overwrite")
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Digest ") {
			nonce := GenerateNonce()
			nonceStore.Store(nonce, time.Now().Add(nonceValidity))
			w.Header().Set("WWW-Authenticate", `Digest realm="`+realmName+`", nonce="`+nonce+`", algorithm="MD5", qop="auth", opaque="`+opaqueValue+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		digestFields := ParseDigestHeader(authHeader[7:])
		username, uri, nonce, response := digestFields["username"], digestFields["uri"], digestFields["nonce"], digestFields["response"]
		// Extract qop-related parameters for RFC 2617 support
		qop, nc, cnonce := digestFields["qop"], digestFields["nc"], digestFields["cnonce"]

		if expiration, ok := nonceStore.Load(nonce); !ok || time.Now().After(expiration.(time.Time)) {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		dbConn, err := db.Open()
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		defer dbConn.Close()

		// Validate using either RFC 2069 or RFC 2617 based on presence of qop
		valid := false
		if qop == "" {
			// RFC 2069 validation
			valid = ValidateDigest(username, uri, r.Method, response, nonce, dbConn)
		} else {
			// RFC 2617 validation with qop
			user, err := models.GetUserByEmail(dbConn, username)
			if err == nil && user != nil {
				password, err := utils.DecodePassword(user.Password)
				if err == nil {
					ha1 := hash([]byte(username + ":" + realmName + ":" + password))
					ha2 := hash([]byte(r.Method + ":" + uri))

					// RFC 2617 with qop
					expectedResponse := hash([]byte(ha1 + ":" + nonce + ":" + nc + ":" + cnonce + ":" + qop + ":" + ha2))
					valid = strings.ToLower(expectedResponse) == strings.ToLower(response)
				}
			}
		}

		if !valid {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		currentUser := extractUserFromAuthHeader(digestFields)
		if currentUser.UserId == "" {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "user", currentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
