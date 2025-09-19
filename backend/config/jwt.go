package config

type JWTConfig struct {
	SecretKey string
}

var jwtConfig *JWTConfig

// LoadJWTConfig loads JWT configuration from environment variables
func LoadJWTConfig() *JWTConfig {
	if jwtConfig == nil {
		jwtConfig = &JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your_secret_key"),
		}
	}
	return jwtConfig
}

// GetJWTConfig returns the loaded JWT configuration
func GetJWTConfig() *JWTConfig {
	if jwtConfig == nil {
		return LoadJWTConfig()
	}
	return jwtConfig
}