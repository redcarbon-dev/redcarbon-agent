package build

// Version is dynamically set by the toolchain
var Version = "DEV"

// DefaultHost is the default host used by the agent
var DefaultHost = "localhost:50051"

// Architecture is dynamically set by the toolchain, and it contains the platform and arch (e.g. linux_amd64)
var Architecture = "dev_arm64"
