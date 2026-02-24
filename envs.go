package utils

import "os"

// EnvOr returns the value of the environment variable named by key,
// or fallback if the variable is not set.
func EnvOr(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func EnvCert(fallback string) string {
	return EnvOr("CERT", fallback)
}

func EnvKey(fallback string) string {
	return EnvOr("KEY", fallback)
}

func EnvHost(fallback string) string {
	return EnvOr("HOST", fallback)
}

func EnvPort(fallback string) string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return fallback
}

func EnvHTTPPort(fallback string) string {
	if port, ok := os.LookupEnv("HTTP_PORT"); ok {
		return ":" + port
	}
	return fallback
}

func EnvUDPPort(fallback string) string {
	if port, ok := os.LookupEnv("UDP_PORT"); ok {
		return ":" + port
	}
	return fallback
}

func EnvAltSvc(fallback string) string {
	return EnvOr("ALT_SVC", fallback)
}

func EnvUI(fallback string) string {
	return EnvOr("UI", fallback)
}

func LookupEnv(key string) *string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return nil
	}
	return &v
}

func LookupEnvPort(key string) *string {
	v := LookupEnv(key)
	if v == nil {
		return nil
	}
	p := ":" + *v
	return &p
}
