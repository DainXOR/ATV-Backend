package types

import (
	"fmt"
	"strconv"
)

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// Shorthand for version V("0.0.0")
func V0() Version {
	return Version{Major: 0, Minor: 0, Patch: 0}
}

func VersionFrom(version string) (Version, error) {
	if len(version) == 1 {
		// If the version string is just a single digit, assume it's the major version
		major, err := strconv.ParseUint(version, 10, 8)
		return Version{Major: uint8(major)}, err
	}

	var major, minor, patch uint8
	_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

// Shorthand for creating a version from a string.
// Only use when you are sure of the input.
// If you are not sure, use VersionFrom instead.
func V(version string) Version {
	if len(version) == 1 {
		// If the version string is just a single digit, assume it's the major version
		major, _ := strconv.ParseUint(version, 10, 8)
		return Version{Major: uint8(major)}
	}

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
func (v Version) GreaterThan(other Version) bool {
	return v.Compare(other) > 0
}
func (v Version) GreaterEq(other Version) bool {
	return v.Compare(other) >= 0
}
func (v Version) LessThan(other Version) bool {
	return v.Compare(other) < 0
}
func (v Version) LessEq(other Version) bool {
	return v.Compare(other) <= 0
}
func (v Version) Equal(other Version) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}

func (v Version) Plus(other Version, others ...Version) Version {
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
func (v Version) Minus(other Version, others ...Version) Version {
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
