package cmd

const (
	pbOwner      = "protocolbuffers"
	pbRepo       = "protobuf"
	pbName       = "pbvm"
	pbDateFormat = "2006.01.02"
)

// pbVersion will be set by goreleaser (see .goreleaser.yml)
var pbVersion = "dev"

// pbCommit will be set by goreleaser (see .goreleaser.yml)
var pbCommit = "unknown"

// pbBuildDt will be set by goreleaser (see .goreleaser.yml)
var pbBuildDt = "now"
