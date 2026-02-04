package config

const (
	Version = "1.1.2"
)

func GetVersion() string {
	return Version
}

func GetFullVersion() string {
	return "Envy v" + Version
}
