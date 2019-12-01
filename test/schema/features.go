package schema

func FeaturesSatisfiesFilter(collection []*DistinguishingFeature, filter FeatureFilterInput) bool {
	return Exist(collection, func(element DistinguishingFeature) bool {
		if element.Name == filter.FeatureName {
			return true
		}
		return false
	})
}

// Exist returns if the collection contains the element that satisfy passed function
func Exist(collection []*DistinguishingFeature, f func(element DistinguishingFeature) bool) bool {
	for _, e := range collection {
		if e == nil {
			continue
		}

		if f(*e) {
			return true
		}
	}
	return false
}
