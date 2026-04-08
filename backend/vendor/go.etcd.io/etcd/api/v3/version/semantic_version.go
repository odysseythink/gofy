package version

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PreRelease string

type SemanticVersion struct {
	Major      int64
	Minor      int64
	Patch      int64
	PreRelease PreRelease
	Metadata   string
}

func splitOff(input *string, delim string) (val string) {
	parts := strings.SplitN(*input, delim, 2)

	if len(parts) == 2 {
		*input = parts[0]
		val = parts[1]
	}

	return val
}

func NewSemanticVersion(version string) *SemanticVersion {
	return MustSemanticVersion(NewSemanticVersionE(version))
}

func NewSemanticVersionE(version string) (*SemanticVersion, error) {
	v := SemanticVersion{}

	if err := v.Set(version); err != nil {
		return nil, err
	}

	return &v, nil
}

// MustSemanticVersion is a helper for wrapping NewVersion and will panic if err is not nil.
func MustSemanticVersion(v *SemanticVersion, err error) *SemanticVersion {
	if err != nil {
		panic(err)
	}
	return v
}

// Set parses and updates v from the given version string. Implements flag.Value
func (v *SemanticVersion) Set(version string) error {
	metadata := splitOff(&version, "+")
	preRelease := PreRelease(splitOff(&version, "-"))
	dotParts := strings.SplitN(version, ".", 3)

	if len(dotParts) != 3 {
		return fmt.Errorf("%s is not in dotted-tri format", version)
	}

	if err := validateIdentifier(string(preRelease)); err != nil {
		return fmt.Errorf("failed to validate pre-release: %v", err)
	}

	if err := validateIdentifier(metadata); err != nil {
		return fmt.Errorf("failed to validate metadata: %v", err)
	}

	parsed := make([]int64, 3)

	for i, v := range dotParts[:3] {
		val, err := strconv.ParseInt(v, 10, 64)
		parsed[i] = val
		if err != nil {
			return err
		}
	}

	v.Metadata = metadata
	v.PreRelease = preRelease
	v.Major = parsed[0]
	v.Minor = parsed[1]
	v.Patch = parsed[2]
	return nil
}

func (v SemanticVersion) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.PreRelease != "" {
		fmt.Fprintf(&buffer, "-%s", v.PreRelease)
	}

	if v.Metadata != "" {
		fmt.Fprintf(&buffer, "+%s", v.Metadata)
	}

	return buffer.String()
}

func (v *SemanticVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}
	return v.Set(data)
}

func (v SemanticVersion) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}

func (v *SemanticVersion) UnmarshalJSON(data []byte) error {
	l := len(data)
	if l == 0 || string(data) == `""` {
		return nil
	}
	if l < 2 || data[0] != '"' || data[l-1] != '"' {
		return errors.New("invalid semver string")
	}
	return v.Set(string(data[1 : l-1]))
}

// Compare tests if v is less than, equal to, or greater than versionB,
// returning -1, 0, or +1 respectively.
func (v SemanticVersion) Compare(versionB SemanticVersion) int {
	if cmp := recursiveCompare(v.Slice(), versionB.Slice()); cmp != 0 {
		return cmp
	}
	return preReleaseCompare(v, versionB)
}

// Equal tests if v is equal to versionB.
func (v SemanticVersion) Equal(versionB SemanticVersion) bool {
	return v.Compare(versionB) == 0
}

// LessThan tests if v is less than versionB.
func (v SemanticVersion) LessThan(versionB SemanticVersion) bool {
	return v.Compare(versionB) < 0
}

// Slice converts the comparable parts of the semver into a slice of integers.
func (v SemanticVersion) Slice() []int64 {
	return []int64{v.Major, v.Minor, v.Patch}
}

func (p PreRelease) Slice() []string {
	preRelease := string(p)
	return strings.Split(preRelease, ".")
}

func preReleaseCompare(versionA SemanticVersion, versionB SemanticVersion) int {
	a := versionA.PreRelease
	b := versionB.PreRelease

	/* Handle the case where if two versions are otherwise equal it is the
	 * one without a PreRelease that is greater */
	if len(a) == 0 && (len(b) > 0) {
		return 1
	} else if len(b) == 0 && (len(a) > 0) {
		return -1
	}

	// If there is a prerelease, check and compare each part.
	return recursivePreReleaseCompare(a.Slice(), b.Slice())
}

func recursiveCompare(versionA []int64, versionB []int64) int {
	if len(versionA) == 0 {
		return 0
	}

	a := versionA[0]
	b := versionB[0]

	if a > b {
		return 1
	} else if a < b {
		return -1
	}

	return recursiveCompare(versionA[1:], versionB[1:])
}

func recursivePreReleaseCompare(versionA []string, versionB []string) int {
	// A larger set of pre-release fields has a higher precedence than a smaller set,
	// if all of the preceding identifiers are equal.
	if len(versionA) == 0 {
		if len(versionB) > 0 {
			return -1
		}
		return 0
	} else if len(versionB) == 0 {
		// We're longer than versionB so return 1.
		return 1
	}

	a := versionA[0]
	b := versionB[0]

	aInt := false
	bInt := false

	aI, err := strconv.Atoi(versionA[0])
	if err == nil {
		aInt = true
	}

	bI, err := strconv.Atoi(versionB[0])
	if err == nil {
		bInt = true
	}

	// Numeric identifiers always have lower precedence than non-numeric identifiers.
	if aInt && !bInt {
		return -1
	} else if !aInt && bInt {
		return 1
	}

	// Handle Integer Comparison
	if aInt && bInt {
		if aI > bI {
			return 1
		} else if aI < bI {
			return -1
		}
	}

	// Handle String Comparison
	if a > b {
		return 1
	} else if a < b {
		return -1
	}

	return recursivePreReleaseCompare(versionA[1:], versionB[1:])
}

// BumpMajor increments the Major field by 1 and resets all other fields to their default values
func (v *SemanticVersion) BumpMajor() {
	v.Major += 1
	v.Minor = 0
	v.Patch = 0
	v.PreRelease = PreRelease("")
	v.Metadata = ""
}

// BumpMinor increments the Minor field by 1 and resets all other fields to their default values
func (v *SemanticVersion) BumpMinor() {
	v.Minor += 1
	v.Patch = 0
	v.PreRelease = PreRelease("")
	v.Metadata = ""
}

// BumpPatch increments the Patch field by 1 and resets all other fields to their default values
func (v *SemanticVersion) BumpPatch() {
	v.Patch += 1
	v.PreRelease = PreRelease("")
	v.Metadata = ""
}

// validateIdentifier makes sure the provided identifier satisfies semver spec
func validateIdentifier(id string) error {
	if id != "" && !reIdentifier.MatchString(id) {
		return fmt.Errorf("%s is not a valid semver identifier", id)
	}
	return nil
}

// reIdentifier is a regular expression used to check that pre-release and metadata
// identifiers satisfy the spec requirements
var reIdentifier = regexp.MustCompile(`^[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*$`)
