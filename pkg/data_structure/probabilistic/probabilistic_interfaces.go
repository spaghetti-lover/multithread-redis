package probabilistic

// FrequencyEstimator defines the interface for Count-Min Sketch
type FrequencyEstimator interface {
	// IncryBy increases the count of item by increment. Multiple items can be increased with one call.
	// Return array of updated min-counts of each of the provided item in the sketch
	// Error if: invalid arguments, missing key, overflow, or wrong key type.
	IncrBy(item string, value uint64) uint64

	// Query returns the count for one or more items in a sketch.
	// Return the min-counts of each of the provided items in the sketch.
	// Error if: invalid arguments, missing key, or wrong key type.
	Count(item string) uint64
}

type MembershipTester interface {
	Add(item string)
	Exist(entry string) bool
}
