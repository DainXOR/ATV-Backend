package types

import "fmt"

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

var (
	VersionZero = Version{Major: 0, Minor: 0, Patch: 0}
	Version100  = Version{Major: 1, Minor: 0, Patch: 0}
	Version010  = Version{Major: 0, Minor: 1, Patch: 0}
	Version001  = Version{Major: 0, Minor: 0, Patch: 1}
)

func VersionFrom(version string) (Version, error) {
	var major, minor, patch uint8
	_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major, Minor: minor, Patch: patch}, nil
}
func V(version string) Version {
	var major, minor, patch uint8
	fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	return Version{Major: major, Minor: minor, Patch: patch}
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
func (v Version) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

func (v Version) Compare(other Version) int8 {
	if v.Major > other.Major {
		return 1
	}
	if v.Major < other.Major {
		return -1
	}
	if v.Minor > other.Minor {
		return 1
	}
	if v.Minor < other.Minor {
		return -1
	}
	if v.Patch > other.Patch {
		return 1
	}
	if v.Patch < other.Patch {
		return -1
	}
	return 0
}
func (v Version) Greater(other Version) bool {
	return v.Compare(other) > 0
}
func (v Version) Less(other Version) bool {
	return v.Compare(other) < 0
}
func (v Version) Equal(other Version) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}

func (v Version) Add(other Version, others ...Version) Version {
	result := Version{
		Major: v.Major + other.Major,
		Minor: v.Minor + other.Minor,
		Patch: v.Patch + other.Patch,
	}
	for _, o := range others {
		result.Major += o.Major
		result.Minor += o.Minor
		result.Patch += o.Patch
	}
	return result
}
func (v Version) Subtract(other Version, others ...Version) Version {
	result := Version{
		Major: v.Major - other.Major,
		Minor: v.Minor - other.Minor,
		Patch: v.Patch - other.Patch,
	}
	for _, o := range others {
		result.Major -= o.Major
		result.Minor -= o.Minor
		result.Patch -= o.Patch
	}
	return result
}
func (v Version) Times(num uint8) Version {
	return Version{
		Major: v.Major * num,
		Minor: v.Minor * num,
		Patch: v.Patch * num,
	}
}
