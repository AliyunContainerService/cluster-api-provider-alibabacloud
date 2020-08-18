package version

import (
	"fmt"
)

var (
	Raw = "v0.0.0-was-not-built-properly"

	// Version = semver.MustParse(strings.TrimLeft(Raw, "v"))

	String = fmt.Sprintf("ClusterAPIProviderAliCLOUD: %s", Raw)
)
