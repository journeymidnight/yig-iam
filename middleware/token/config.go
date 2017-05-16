package token

const (
	DefaultContextKey = "token"
	DefaultCookieKey = "yig_iam_login_token"
)

// Config is a struct for specifying configuration options for the jwt middleware.
type Config struct {
	ContextKey string
	CookieKey string
}
