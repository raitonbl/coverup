package config

type Manifest struct {
	Services []Service         `yaml:"services"`
	Config   map[string]Config `yaml:"config"`
}

type Service struct {
	Config    *Config            `yaml:"config"`
	QueueURL  string             `yaml:"queue_url,omitempty"`
	Resources map[string]*Config `yaml:"resources,omitempty"`
}

// Common

type Config struct {
	Region      string       `yaml:"region"`
	Credentials *Credentials `yaml:"credentials"`
}

type Credentials struct {
	Type            string           `yaml:"type"`
	StaticProvider  *StaticProvider  `yaml:"static-provider,omitempty"`
	ProfileProvider *ProfileProvider `yaml:"profile-provider,omitempty"`
}

type StaticProvider struct {
	AccessKeyID     string `yaml:"access-key-id"`
	SecretAccessKey string `yaml:"secret-access-key"`
}

type ProfileProvider struct {
	ProfileName string `yaml:"profile-name"`
}
