package app

// Open open
func Open() error {
	for _, f := range resources {
		if e := f(); e != nil {
			return e
		}
	}
	return nil
}

// OpenFunc load viper config
type OpenFunc func() error

var resources []OpenFunc

// RegisterResource register resource
func RegisterResource(args ...OpenFunc) {
	resources = append(resources, args...)
}
