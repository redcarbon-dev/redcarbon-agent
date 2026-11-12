package build

// Version is dynamically set by the toolchain
var Version = "DEV"

// DefaultHost is the default host used by the agent
var DefaultHost = "localhost:50051"

// GitHubAPIURL is the GitHub REST API the agent fetches its releases from
var GitHubAPIURL = "https://api.github.com"

// OutboundIPProbeAddr is the UDP address dialed (no packet is sent) to learn the agent's outbound IP
var OutboundIPProbeAddr = "8.8.8.8:80"

// Architecture is dynamically set by the toolchain, and it contains the platform and arch (e.g. linux_amd64)
var Architecture = "darwin_arm64"

// UserAgent is the User-Agent header the agent sends on outbound HTTP requests.
func UserAgent() string {
	return "redcarbon-agent/" + Version
}
