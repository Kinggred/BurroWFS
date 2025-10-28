package schemas

import "encoding/xml"

type LockInfo struct {
	XMLName   xml.Name
	LockScope LockScope `xml:"lockscope"`
	LockType  LockType  `xml:"locktype"`
	Owner     LockOwner `xml:"owner"`
}

type LockScope struct {
	Exclusive *struct{} `xml:"exclusive,omitempty"`
	Shared    *struct{} `xml:"shared,omitempty"`
}

type LockType struct {
	Write *struct{} `xml:"write,omitempty"`
}

type LockOwner struct {
	Href string `xml:"href"`
}
