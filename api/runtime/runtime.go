package runtime

import "fmt"

// VersionString returns a version, if this is a release version, or a pseudo-version using a git commit hash otherwise.
func (r *Runtime) VersionString() string {
	if r.ReleaseVersion != nil {
		return r.GetReleaseVersion()
	} else if r.CommitTime != nil && r.CommitHash != "" {
		return fmt.Sprintf("v0.0.0-%s-%s", r.CommitTime.AsTime().Format("20060102150405"), r.CommitHash[0:12])
	} else {
		return "v0.0.0"
	}
}
