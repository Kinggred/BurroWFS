package webdav

import "encoding/xml"

type LockInfo struct {
	XMLName xml.Name   `xml:"lockinfo"`
	Scope   LockScope  `xml:"lockscope"`
	Type    LockType   `xml:"locktype"`
	Owner   *LockOwner `xml:"owner,omitempty"`
}

type LockScope struct {
	Exclusive *struct{} `xml:"exclusive,omitempty"`
	Shared    *struct{} `xml:"shared,omitempty"`
}

func (s LockScope) String() string {
	if s.Shared != nil {
		return "shared"
	}
	return "exclusive"
}

type LockType struct {
	Write *struct{} `xml:"write"`
}

type LockOwner struct {
	Href string `xml:"href,omitempty"`
	Raw  string `xml:",innerxml"`
}
