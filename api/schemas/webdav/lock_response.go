package webdav

import (
	"encoding/xml"
	"strings"
)

type nsDAV string

const DAV nsDAV = "DAV:"

type LockTypeXML struct {
	Write struct{} `xml:"D:write"`
}

type LockScopeXML struct {
	Exclusive *struct{} `xml:"D:exclusive,omitempty"`
	Shared    *struct{} `xml:"D:shared,omitempty"`
}

type OwnerXML struct {
	Href string `xml:"D:href,omitempty"`
}

type HrefXML struct {
	Href string `xml:"D:href"`
}

type LockTokenXML struct {
	Href string `xml:"D:href"`
}

type ActiveLockXML struct {
	XMLName   xml.Name     `xml:"D:activelock"`
	LockType  LockTypeXML  `xml:"D:locktype"`
	LockScope LockScopeXML `xml:"D:lockscope"`
	Depth     string       `xml:"D:depth"` // "0" lub "infinity"
	Owner     *OwnerXML    `xml:"D:owner,omitempty"`
	Timeout   string       `xml:"D:timeout,omitempty"` // "Second-600" albo "Infinite"
	LockToken LockTokenXML `xml:"D:locktoken"`
	LockRoot  *HrefXML     `xml:"D:lockroot,omitempty"` // D:href do zablokowanego zasobu
}

type LockDiscoveryXML struct {
	XMLName    xml.Name      `xml:"D:lockdiscovery"`
	ActiveLock ActiveLockXML `xml:"D:activelock"`
}

type LockResponse struct {
	XMLName       xml.Name         `xml:"D:prop"`
	XmlnsD        string           `xml:"xmlns:D,attr"`
	LockDiscovery LockDiscoveryXML `xml:"D:lockdiscovery"`
}

func NewLockScopeXML(shared bool) LockScopeXML {
	if shared {
		return LockScopeXML{Shared: &struct{}{}, Exclusive: nil}
	}
	return LockScopeXML{Exclusive: &struct{}{}, Shared: nil}
}

func NewOwnerXML(ownerHref string) *OwnerXML {
	ownerHref = strings.TrimSpace(ownerHref)
	if ownerHref == "" {
		return nil
	}
	return &OwnerXML{Href: ownerHref}
}

func OwnerFromLockInfo(info LockInfo) *OwnerXML {
	if info.Owner == nil {
		return nil
	}
	owner := strings.TrimSpace(info.Owner.Href)
	if owner == "" {
		owner = strings.TrimSpace(info.Owner.Raw)
	}
	return NewOwnerXML(owner)
}

func ScopeFromLockInfo(info LockInfo) LockScopeXML {
	shared := info.Scope.Shared != nil
	return NewLockScopeXML(shared)
}

// NewActiveLock buduje strukturę ActiveLockXML dla odpowiedzi LOCK.
// depth: "0" lub "infinity"
// timeout: np. "Second-600" lub "Infinite"
// token: pełny format, np. "opaquelocktoken:..."
// lockRootHref: zasób podlegający blokadzie
func NewActiveLock(depth, timeout, token, lockRootHref string, scope LockScopeXML, owner *OwnerXML) ActiveLockXML {
	al := ActiveLockXML{
		LockType:  LockTypeXML{Write: struct{}{}},
		LockScope: scope,
		Depth:     depth,
		Owner:     owner,
		Timeout:   timeout,
		LockToken: LockTokenXML{Href: token},
	}
	lockRootHref = strings.TrimSpace(lockRootHref)
	if lockRootHref != "" {
		al.LockRoot = &HrefXML{Href: lockRootHref}
	}
	return al
}

func NewLockResponse(active ActiveLockXML) LockResponse {
	return LockResponse{
		XmlnsD: string(DAV),
		LockDiscovery: LockDiscoveryXML{
			ActiveLock: active,
		},
	}
}

func BuildLockResponseFromRequest(info LockInfo, depth, timeout, token, lockRootHref string) LockResponse {
	owner := OwnerFromLockInfo(info)
	scope := ScopeFromLockInfo(info)
	active := NewActiveLock(depth, timeout, token, lockRootHref, scope, owner)
	return NewLockResponse(active)
}

func BuildLockResponse(depth, timeout, token, lockRootHref, ownerHref string, shared bool) LockResponse {
	owner := NewOwnerXML(ownerHref)
	scope := NewLockScopeXML(shared)
	active := NewActiveLock(depth, timeout, token, lockRootHref, scope, owner)
	return NewLockResponse(active)
}
