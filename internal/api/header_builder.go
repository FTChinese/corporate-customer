package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"net/http"
)

const (
	keyFtcID         = "X-User-Id"
	keyUnionID       = "X-Union-Id"
	keyWxAppID       = "X-App-Id"
	keyClientType    = "X-Client-Type"
	keyClientVersion = "X-Client-Version"
	keyUserIP        = "X-User-Ip"
	keyUserAgent     = "X-User-Agent"
)

type HeaderBuilder struct {
	h http.Header
}

func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{
		h: http.Header{},
	}
}

func ReaderIDsHeader(c reader.PassportClaims) *HeaderBuilder {
	b := NewHeaderBuilder()
	if c.FtcID != "" {
		b.WithFtcID(c.FtcID)
	}

	if c.UnionID.Valid {
		b.WithUnionID(c.UnionID.String)
	}

	return b
}

func (b *HeaderBuilder) WithPlatformWeb() *HeaderBuilder {
	b.h.Set(keyClientType, "web")
	return b
}

func (b *HeaderBuilder) WithClientVersion(v string) *HeaderBuilder {
	b.h.Set(keyClientVersion, v)
	return b
}

func (b *HeaderBuilder) WithUserIP(ip string) *HeaderBuilder {
	b.h.Set(keyUserIP, ip)
	return b
}

func (b *HeaderBuilder) WithUserAgent(ua string) *HeaderBuilder {
	b.h.Set(keyUserAgent, ua)
	return b
}

func (b *HeaderBuilder) WithWxAppID(id string) *HeaderBuilder {
	b.h.Set(keyWxAppID, id)
	return b
}

func (b *HeaderBuilder) WithFtcID(id string) *HeaderBuilder {
	b.h.Set(keyFtcID, id)
	return b
}

func (b *HeaderBuilder) WithUnionID(id string) *HeaderBuilder {
	b.h.Set(keyUnionID, id)
	return b
}

func (b *HeaderBuilder) Build() http.Header {
	return b.h
}
