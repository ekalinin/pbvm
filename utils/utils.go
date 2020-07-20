package utils

import (
	"fmt"
	"runtime"
	"strings"
)

var archs = map[string]string{
	"386":   "x86_32",
	"amd64": "x86_64",
}

// GetArch returns current arch
func GetArch() string {
	arch, ok := archs[runtime.GOARCH]
	if ok {
		return arch
	}
	return runtime.GOARCH
}

// IsSuitableAsset returns true if asset is suitable for download for current system
func IsSuitableAsset(assetName, arch, os string) bool {
	return strings.HasPrefix(assetName, "protoc") &&
		strings.Contains(assetName, fmt.Sprintf("-%s-", os)) &&
		strings.Contains(assetName, fmt.Sprintf("-%s.", arch))
}
