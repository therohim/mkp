package config

type AuthConfig struct {
	StoreKey     string
	Lifetime     int // In day(s)
	MaxAge       int // In second
	CookieDomain string
	HTTPOnly     bool
}

func NewAuthConfig(c Config) AuthConfig {
	return AuthConfig{
		StoreKey:     "AuthTestMkp",
		Lifetime:     365,    // 1 year
		MaxAge:       604800, // about 7 days
		CookieDomain: c.Get("AUTH_COOKIE_DOMAIN"),
		HTTPOnly:     true,
	}
}
