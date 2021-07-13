package version

import (
<<<<<<< HEAD
<<<<<<< HEAD
	"fmt"
<<<<<<< HEAD
	"runtime"
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 24c35849 (fix stop ecs instance func)
=======
	"fmt"
	"runtime"
>>>>>>> 836a3e90 (update README)
	"strings"

	"github.com/blang/semver"
)

<<<<<<< HEAD
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
=======
var (
	// Raw is the string representation of the version. This will be replaced with the calculated version at build time.
	Raw = "v0.1.0-alpha.1"

	// Version is semver representation of the version.
	Version = semver.MustParse(strings.TrimLeft(Raw, "v"))
<<<<<<< HEAD

	// String is the human-friendly representation of the version.
<<<<<<< HEAD
	String = fmt.Sprintf("ClusterAPIProvideralibabacloud %s", Raw)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	String = fmt.Sprintf("ClusterAPIProviderAlibabaCloud %s", Raw)
>>>>>>> ecfeb08f (remove unused code)
=======
>>>>>>> 24c35849 (fix stop ecs instance func)
)

<<<<<<< HEAD
=======
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

>>>>>>> 836a3e90 (update README)
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

<<<<<<< HEAD
<<<<<<< HEAD
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
=======
func printShortDirtyVersionInfo() {
	fmt.Printf("Version Info: GitReleaseTag: %q, MajorVersion: %q, MinorVersion:%q, GitReleaseCommit:%q, GitTreeState:%q\n",
=======
func PrintShortDirtyVersionInfo() string {
	return fmt.Sprintf("Version Info: GitReleaseTag: %q, MajorVersion: %q, MinorVersion:%q, GitReleaseCommit:%q, GitTreeState:%q",
>>>>>>> 6a93b4ce (print version)
		gitReleaseTag, gitMajor, gitMinor, gitReleaseCommit, gitTreeState)
}

func PrintShortCleanVersionInfo() string {
	return fmt.Sprintf("Version Info: GitReleaseTag: %q, MajorVersion: %q, MinorVersion:%q", gitReleaseTag, gitMajor, gitMinor)
}

<<<<<<< HEAD
func printVerboseVersionInfo() {
	fmt.Println("Version Info:")
	fmt.Printf("GitReleaseTag: %q, Major: %q, Minor: %q, GitRelaseCommit: %q\n", gitReleaseTag, gitMajor, gitMinor, gitReleaseCommit)
	fmt.Printf("Git Branch: %q\n", gitBranch)
	fmt.Printf("Git commit: %q\n", gitCommit)
	fmt.Printf("Git tree state: %q\n", gitTreeState)
>>>>>>> 836a3e90 (update README)
=======
func PrintVerboseVersionInfo() string {
	return fmt.Sprintf("Version Info: GitReleaseTag: %q, Major: %q, Minor: %q, GitRelaseCommit: %q,Git Branch: %q,Git commit: %q,Git tree state: %q",
		gitReleaseTag, gitMajor, gitMinor, gitReleaseCommit, gitBranch, gitCommit, gitTreeState)
>>>>>>> 6a93b4ce (print version)
}
