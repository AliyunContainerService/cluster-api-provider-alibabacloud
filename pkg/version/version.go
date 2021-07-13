package version

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/blang/semver"
)

var (
	// Raw is the string representation of the version. This will be replaced with the calculated version at build time.
	Raw = "v0.1.0-alpha.1"

	// Version is semver representation of the version.
	Version = semver.MustParse(strings.TrimLeft(Raw, "v"))
)

var (
	gitMajor         string // major version, always numeric
	gitMinor         string // minor version, numeric possibly followed by "+"
	gitVersion       string // semantic version, derived by build scripts
	gitCommit        string // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState     string // state of git tree, either "clean" or "dirty"
	buildDate        string // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	gitReleaseTag    string // gitReleaseTag is the git tag from which this binary is released
	gitReleaseCommit string // gitReleaseCommit is the commit corresponding to the GitReleaseTag
	gitBranch        string // gitBranch is the branch from which this binary was built
)

type Info struct {
	Major            string `json:"major,omitempty"`
	Minor            string `json:"minor,omitempty"`
	GitVersion       string `json:"gitVersion,omitempty"`
	GitCommit        string `json:"gitCommit,omitempty"`
	GitTreeState     string `json:"gitTreeState,omitempty"`
	BuildDate        string `json:"buildDate,omitempty"`
	GoVersion        string `json:"goVersion,omitempty"`
	Compiler         string `json:"compiler,omitempty"`
	Platform         string `json:"platform,omitempty"`
	GitReleaseTag    string `json:"gitReleaseTag,omitempty"`
	GitReleaseCommit string `json:"gitReleaseCommit,omitempty"`
	GitBranch        string `json:"gitBranch,omitempty"`
}

func GetVersionInfo() Info {
	return Info{
		Major:            gitMajor,
		Minor:            gitMinor,
		GitVersion:       gitVersion,
		GitCommit:        gitCommit,
		GitTreeState:     gitTreeState,
		BuildDate:        buildDate,
		GoVersion:        runtime.Version(),
		Compiler:         runtime.Compiler,
		Platform:         fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GitReleaseTag:    gitReleaseTag,
		GitReleaseCommit: gitReleaseCommit,
		GitBranch:        gitBranch,
	}
}

func (info Info) String() string {
	return info.GitVersion
}

func PrintShortDirtyVersionInfo() string {
	return fmt.Sprintf("Version Info: GitReleaseTag: %q, MajorVersion: %q, MinorVersion:%q, GitReleaseCommit:%q, GitTreeState:%q",
		gitReleaseTag, gitMajor, gitMinor, gitReleaseCommit, gitTreeState)
}

func PrintShortCleanVersionInfo() string {
	return fmt.Sprintf("Version Info: GitReleaseTag: %q, MajorVersion: %q, MinorVersion:%q", gitReleaseTag, gitMajor, gitMinor)
}

func PrintVerboseVersionInfo() string {
	return fmt.Sprintf("Version Info: GitReleaseTag: %q, Major: %q, Minor: %q, GitRelaseCommit: %q,Git Branch: %q,Git commit: %q,Git tree state: %q",
		gitReleaseTag, gitMajor, gitMinor, gitReleaseCommit, gitBranch, gitCommit, gitTreeState)
}
