package aoxnet

// Feature identifies an optional AOXNet capability negotiated during handshake.
type Feature string

const (
	FeatureRequestID    Feature = "request-id"
	FeatureCRC32        Feature = "crc32"
	FeaturePing         Feature = "ping"
	FeatureTLSTransport Feature = "tls-transport"
	FeatureMultiplexing Feature = "multiplexing"
	FeatureFlowControl  Feature = "flow-control"
	FeatureGoAway       Feature = "goaway"
)

// DefaultFeatures returns the capabilities implemented by this AOXNet runtime.
func DefaultFeatures() []string {
	return []string{
		string(FeatureRequestID),
		string(FeatureCRC32),
		string(FeaturePing),
		string(FeatureTLSTransport),
		string(FeatureMultiplexing),
		string(FeatureFlowControl),
		string(FeatureGoAway),
	}
}

// NegotiateFeatures returns the ordered intersection of local and remote capabilities.
func NegotiateFeatures(local, remote []string) []string {
	result := make([]string, 0, len(local))
	for _, feature := range local {
		if containsString(remote, feature) && !containsString(result, feature) {
			result = append(result, feature)
		}
	}
	return result
}

// HasFeature reports whether a negotiated feature is present.
func HasFeature(features []string, feature Feature) bool {
	return containsString(features, string(feature))
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
