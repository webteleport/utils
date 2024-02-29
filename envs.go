package utils

import "os"

func EnvCert(fallback string) string {
	if cert, ok := os.LookupEnv("CERT"); ok {
		return cert
	}
	return fallback
}

func EnvKey(fallback string) string {
	if key, ok := os.LookupEnv("KEY"); ok {
		return key
	}
	return fallback
}

func EnvHost(fallback string) string {
	if host, ok := os.LookupEnv("HOST"); ok {
		return host
	}
	return fallback
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
	if altsvc, ok := os.LookupEnv("ALT_SVC"); ok {
		return altsvc
	}
	return fallback
}

func EnvUI(fallback string) string {
	if ui, ok := os.LookupEnv("UI"); ok {
		return ui
	}
	return fallback
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
