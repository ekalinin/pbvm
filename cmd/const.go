package cmd

const (
	pbOwner      = "protocolbuffers"
	pbRepo       = "protobuf"
	pbName       = "pbvm"
	pbDateFormat = "2006.01.02"
)

// Version will be set by goreleaser (see .goreleaser.yml)
var pbVersion string

// Commit will be set by goreleaser (see .goreleaser.yml)
var pbCommit string
