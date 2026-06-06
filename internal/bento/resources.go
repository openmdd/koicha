package bento

// ResourceView describes the Kafka resources visible by default in a bento.
//
// Topics and ConsumerGroups are filtered independently using the same
// ResourceScope mechanism.
//
// When AllowOutOfScope is true the user may navigate to any resource in the
// cluster even if it does not match the defined scopes; the scopes only set the
// default view in that case.
type ResourceView struct {
	Topics          ResourceScope `yaml:"topics"`
	ConsumerGroups  ResourceScope `yaml:"consumerGroups"`
	AllowOutOfScope bool          `yaml:"allowOutOfScope"`
}

// ResourceScope limits a Kafka resource list.
//
// Matching order: Include is evaluated first. An empty Include matches all
// resources. Exclude is then applied to the result. A resource must match at
// least one Include pattern (or Include must be empty) AND must not match any
// Exclude pattern to be visible.
type ResourceScope struct {
	Include []ResourcePattern `yaml:"include,omitempty"`
	Exclude []ResourcePattern `yaml:"exclude,omitempty"`
}

// ResourcePatternKind controls how a ResourcePattern value is interpreted.
type ResourcePatternKind string

const (
	// ResourcePatternExact matches a resource whose name equals Value exactly.
	ResourcePatternExact ResourcePatternKind = "exact"
	// ResourcePatternPrefix matches a resource whose name starts with Value.
	ResourcePatternPrefix ResourcePatternKind = "prefix"
	// ResourcePatternRegex matches a resource whose name satisfies the RE2
	// regular expression in Value.
	ResourcePatternRegex ResourcePatternKind = "regex"
)

// ResourcePattern is a single include or exclude rule for a ResourceScope.
type ResourcePattern struct {
	Kind  ResourcePatternKind `yaml:"kind"`
	Value string              `yaml:"value"`
}
