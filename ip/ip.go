package ip

import (
	"net/http"
	"strings"
)

func GetRemoteAddr(r *http.Request) string {
	ip := removePort(r.RemoteAddr)
	if len(r.Header.Get("X-FORWARDED-FOR")) > 0 {
		ip = r.Header.Get("X-FORWARDED-FOR")
	}
	return ip
}

func removePort(ip string) string {
	split := strings.Split(ip, ":")
	if len(split) < 2 {
		return ip
	}
	return strings.Join(split[:len(split)-1], ":")
}
