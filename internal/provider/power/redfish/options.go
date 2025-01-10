package redfish

// Options is a struct that holds the RedFish configuration options.
type Options struct {
	UseAlways                 bool
	UseWhenAvailable          bool
	UseHTTPS                  bool
	InsecureSkipTLSVerify     bool
	SetBootSourceOverrideMode bool
	Port                      int
}

// DefaultOptions is the default RedFish configuration options.
var DefaultOptions = Options{
	UseAlways:                 false,
	UseWhenAvailable:          true,
	UseHTTPS:                  true,
	InsecureSkipTLSVerify:     true,
	Port:                      443,
	SetBootSourceOverrideMode: true,
}
