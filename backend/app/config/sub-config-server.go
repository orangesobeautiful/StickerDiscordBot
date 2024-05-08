package config

import (
	"net/http"
	"time"
)

type Server interface {
	GetAddr() string
	GetImgURL() string
	GetSessionKey() ServerSessionKey
	GetCookie() Cookie
	GetCORS() CORS
}

type ServerSessionKey interface {
	GetUserAuth() SessionKey
}

type SessionKey interface {
	GetAuthentication() []byte
	GetEncryption() []byte
	SessionKeyPair() [][]byte
}

type Cookie interface {
	GetMaxAge() time.Duration
	GetSecure() bool
	GetHTTPOnly() bool
	GetSameSite() http.SameSite
}

type CORS interface {
	GetAllowAllOrigins() bool
	GetAllowOrigins() []string
	GetAllowMethods() []string
	GetAllowHeaders() []string
	GetExposeHeaders() []string
	GetAllowCredentials() bool
	GetMaxAge() time.Duration
}

var _ Server = (*server)(nil)

type server struct {
	Addr       string
	ImgURL     string
	SessionKey *serverSessionKey
	Cookie     *cookie
	CORS       *cors
}

func (s *server) GetAddr() string {
	return s.Addr
}

func (s *server) GetImgURL() string {
	return s.ImgURL
}

func (s *server) GetSessionKey() ServerSessionKey {
	return s.SessionKey
}

func (s *server) GetCookie() Cookie {
	return s.Cookie
}

func (s *server) GetCORS() CORS {
	return s.CORS
}

var _ ServerSessionKey = (*serverSessionKey)(nil)

type serverSessionKey struct {
	UserAuth *sessionKey
}

func (s *serverSessionKey) GetUserAuth() SessionKey {
	return s.UserAuth
}

var _ SessionKey = (*sessionKey)(nil)

type sessionKey struct {
	Authentication []byte
	Encryption     []byte
}

func (s *sessionKey) GetAuthentication() []byte {
	return s.Authentication
}

func (s *sessionKey) GetEncryption() []byte {
	return s.Encryption
}

func (s *sessionKey) SessionKeyPair() [][]byte {
	keyPair := [][]byte{s.Authentication}

	if len(s.Encryption) != 0 {
		keyPair = append(keyPair, s.Encryption)
	}
	return keyPair
}

var _ Cookie = (*cookie)(nil)

type cookie struct {
	MaxAge   time.Duration
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

func (c *cookie) GetMaxAge() time.Duration {
	return c.MaxAge
}

func (c *cookie) GetSecure() bool {
	return c.Secure
}

func (c *cookie) GetHTTPOnly() bool {
	return c.HTTPOnly
}

func (c *cookie) GetSameSite() http.SameSite {
	return c.SameSite
}

var _ CORS = (*cors)(nil)

type cors struct {
	AllowAllOrigins  bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

func (c *cors) GetAllowAllOrigins() bool {
	return c.AllowAllOrigins
}

func (c *cors) GetAllowOrigins() []string {
	return c.AllowOrigins
}

func (c *cors) GetAllowMethods() []string {
	return c.AllowMethods
}

func (c *cors) GetAllowHeaders() []string {
	return c.AllowHeaders
}

func (c *cors) GetExposeHeaders() []string {
	return c.ExposeHeaders
}

func (c *cors) GetAllowCredentials() bool {
	return c.AllowCredentials
}

func (c *cors) GetMaxAge() time.Duration {
	return c.MaxAge
}
