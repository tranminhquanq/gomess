package utils

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/tranminhquanq/gomess/internal/config"
)

// GetIPAddress returns the real IP address of the HTTP request. It parses the
// X-Forwarded-For header.
func GetIPAddress(r *http.Request) string {
	if r.Header != nil {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		if xForwardedFor != "" {
			ips := strings.Split(xForwardedFor, ",")
			for i := range ips {
				ips[i] = strings.TrimSpace(ips[i])
			}

			for _, ip := range ips {
				if ip != "" {
					parsed := net.ParseIP(ip)
					if parsed == nil {
						continue
					}

					return parsed.String()
				}
			}
		}
	}

	ipPort := r.RemoteAddr
	ip, _, err := net.SplitHostPort(ipPort)
	if err != nil {
		return ipPort
	}

	return ip
}

// getRedirectTo tries extract redirect url from header or from query params
func getRedirectTo(r *http.Request) (reqref string) {
	reqref = r.Header.Get("redirect_to")
	if reqref != "" {
		return
	}

	if err := r.ParseForm(); err == nil {
		reqref = r.Form.Get("redirect_to")
	}

	return
}

func IsRedirectURLValid(config *config.GlobalConfiguration, redirectURL string) bool {
	if redirectURL == "" {
		return false
	}

	base, berr := url.Parse(config.SiteURL)
	refurl, rerr := url.Parse(redirectURL)

	// As long as the referrer came from the site, we will redirect back there
	if berr == nil && rerr == nil && base.Hostname() == refurl.Hostname() {
		return true
	}

	// For case when user came from mobile app or other permitted resource - redirect back
	for _, pattern := range config.URIAllowListMap {
		if pattern.Match(redirectURL) {
			return true
		}
	}

	return false
}

func GetReferrer(r *http.Request, config *config.GlobalConfiguration) string {
	// try get redirect url from query or post data first
	reqref := getRedirectTo(r)
	if IsRedirectURLValid(config, reqref) {
		return reqref
	}

	// instead try referrer header value
	reqref = r.Referer()
	if IsRedirectURLValid(config, reqref) {
		return reqref
	}

	return config.SiteURL
}
