package mapper

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/config"
)

// IPResolver handles hostname to IP resolution with caching
type IPResolver struct {
	cache    map[string]cacheEntry
	mu       sync.RWMutex
	cacheTTL time.Duration
}

type cacheEntry struct {
	ip        string
	expiresAt time.Time
}

// Default TTL from environment or 5 minutes
var defaultCacheTTL = time.Duration(config.GetEnvInt("IP_RESOLVER_CACHE_TTL_SECONDS", 300)) * time.Second
var defaultResolver = NewIPResolver(defaultCacheTTL)

// NewIPResolver creates a new IP resolver with cache TTL
func NewIPResolver(ttl time.Duration) *IPResolver {
	return &IPResolver{
		cache:    make(map[string]cacheEntry),
		cacheTTL: ttl,
	}
}

// ResolveIP resolves a hostname to an IP address with caching
// If the input is already an IP, returns it directly
// Falls back to "0.0.0.0" if resolution fails
func (r *IPResolver) ResolveIP(hostOrIP string) string {
	if hostOrIP == "" {
		return "0.0.0.0"
	}

	// Clean up the input (might have port attached)
	host := hostOrIP
	if colonIdx := strings.LastIndex(hostOrIP, ":"); colonIdx != -1 {
		// Check if it's IPv6 or host:port
		if !strings.Contains(hostOrIP, "[") {
			host = hostOrIP[:colonIdx]
		}
	}

	// If it's already a valid IP, return it
	if ip := net.ParseIP(host); ip != nil {
		return ip.String()
	}

	// Check cache first
	r.mu.RLock()
	entry, exists := r.cache[host]
	r.mu.RUnlock()

	if exists && time.Now().Before(entry.expiresAt) {
		return entry.ip
	}

	// Resolve the hostname
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		// Cache the failure too to avoid repeated lookups
		r.cacheResult(host, "0.0.0.0")
		return "0.0.0.0"
	}

	// Prefer IPv4 addresses
	resultIP := ips[0].String()
	for _, ip := range ips {
		if ip.To4() != nil {
			resultIP = ip.String()
			break
		}
	}

	r.cacheResult(host, resultIP)
	return resultIP
}

func (r *IPResolver) cacheResult(host, ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[host] = cacheEntry{
		ip:        ip,
		expiresAt: time.Now().Add(r.cacheTTL),
	}
}

// ResolveHostIP is a convenience function using the default resolver
func ResolveHostIP(hostOrIP string) string {
	return defaultResolver.ResolveIP(hostOrIP)
}
