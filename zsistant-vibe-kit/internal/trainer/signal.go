package trainer

// Dimension tracks a style preference with a score.
type Dimension struct {
	Name   string
	Value  string
	Score  int // accumulated evidence weight
}

// Signal is a detected user preference from a single interaction.
type Signal struct {
	Dimension string
	Value     string
	Evidence  string // the snippet that triggered this signal
}

// StyleProfile aggregates detected signals into scored dimensions.
type StyleProfile struct {
	Dimensions map[string]*Dimension
}

// NewStyleProfile creates an empty profile.
func NewStyleProfile() *StyleProfile {
	return &StyleProfile{
		Dimensions: make(map[string]*Dimension),
	}
}

// Apply adds a signal to the profile, accumulating scores.
func (sp *StyleProfile) Apply(s Signal) {
	key := s.Dimension + ":" + s.Value
	if d, ok := sp.Dimensions[key]; ok {
		d.Score++
	} else {
		sp.Dimensions[key] = &Dimension{
			Name:  s.Dimension,
			Value: s.Value,
			Score: 1,
		}
	}
}

// Top returns the highest-scoring value for each dimension.
func (sp *StyleProfile) Top() map[string]*Dimension {
	best := make(map[string]*Dimension)
	for _, d := range sp.Dimensions {
		current, ok := best[d.Name]
		if !ok || d.Score > current.Score {
			best[d.Name] = d
		}
	}
	return best
}

// Confident returns true if the top value for a dimension exceeds threshold.
func (sp *StyleProfile) Confident(dimension string, threshold int) bool {
	top := sp.Top()
	d, ok := top[dimension]
	return ok && d.Score >= threshold
}
