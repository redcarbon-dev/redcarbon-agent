package build

// Version is dynamically set by the toolchain
var Version = "DEV"

// DefaultHost is the default host used by the agent
var DefaultHost = "localhost:50051"

// Architecture is dynamically set by the toolchain, and it contains the platform and arch (e.g. linux_amd64)
var Architecture = "darwin_arm64"

// UserAgent is the User-Agent header the agent sends on outbound HTTP requests.
func UserAgent() string {
	return "redcarbon-agent/" + Version
}
