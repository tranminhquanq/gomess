package config

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobwas/glob"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Time is used to represent timestamps in the configuration, as envconfig has
// trouble parsing empty strings, due to time.Time.UnmarshalText().
type Time struct {
	time.Time
}

func (t *Time) UnmarshalText(text []byte) error {
	trimed := bytes.TrimSpace(text)

	if len(trimed) < 1 {
		t.Time = time.Time{}
	} else {
		if err := t.Time.UnmarshalText(trimed); err != nil {
			return err
		}
	}

	return nil
}

type APIConfiguration struct {
	ID                 string `json:"id" envconfig:"API_ID" required:"true"`
	Host               string
	Port               string `envconfig:"PORT" default:"9999"`
	Endpoint           string
	RequestIDHeader    string        `envconfig:"REQUEST_ID_HEADER"`
	ExternalURL        string        `json:"external_url" envconfig:"API_EXTERNAL_URL" required:"true"`
	MaxRequestDuration time.Duration `json:"max_request_duration" split_words:"true" default:"10s"`
}

func (a *APIConfiguration) Validate() error {
	_, err := url.ParseRequestURI(a.ExternalURL)
	if err != nil {
		return err
	}

	return nil
}

type CORSConfiguration struct {
	AllowedHeaders []string `json:"allowed_headers" split_words:"true"`
}

func (c *CORSConfiguration) AllAllowedHeaders(defaults []string) []string {
	set := make(map[string]bool)
	for _, header := range defaults {
		set[header] = true
	}

	var result []string
	result = append(result, defaults...)

	for _, header := range c.AllowedHeaders {
		if !set[header] {
			result = append(result, header)
		}

		set[header] = true
	}

	return result
}

// DBConfiguration holds all the database related configuration.
type DBConfiguration struct {
	Driver    string `json:"driver" required:"true"`
	URL       string `json:"url" envconfig:"DATABASE_URL" required:"true"`
	Namespace string `json:"namespace" envconfig:"DB_NAMESPACE" default:"auth"`
	// MaxPoolSize defaults to 0 (unlimited).
	MaxPoolSize       int           `json:"max_pool_size" split_words:"true"`
	MaxIdlePoolSize   int           `json:"max_idle_pool_size" split_words:"true"`
	ConnMaxLifetime   time.Duration `json:"conn_max_lifetime,omitempty" split_words:"true"`
	ConnMaxIdleTime   time.Duration `json:"conn_max_idle_time,omitempty" split_words:"true"`
	HealthCheckPeriod time.Duration `json:"health_check_period" split_words:"true"`
	MigrationsPath    string        `json:"migrations_path" split_words:"true" default:"./migrations"`
	CleanupEnabled    bool          `json:"cleanup_enabled" split_words:"true" default:"false"`
}

func (c *DBConfiguration) Validate() error {
	return nil
}

// GlobalConfiguration holds all the configuration that applies to all instances.
type GlobalConfiguration struct {
	API     APIConfiguration
	CORS    CORSConfiguration
	DB      DBConfiguration
	Tracing TracingConfig
	Metrics MetricsConfig

	SiteURL         string   `json:"site_url" split_words:"true" required:"true"`
	URIAllowList    []string `json:"uri_allow_list" split_words:"true"`
	URIAllowListMap map[string]glob.Glob
}

// ApplyDefaults sets defaults for a GlobalConfiguration
func (config *GlobalConfiguration) ApplyDefaults() error {
	return nil
}

// Validate validates all of configuration.
func (c *GlobalConfiguration) Validate() error {
	validatables := []interface {
		Validate() error
	}{
		&c.API,
		&c.DB,
		&c.Tracing,
		&c.Metrics,
	}

	for _, validatable := range validatables {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// LoadGlobalFromEnv will return a new *GlobalConfiguration value from the
// currently configured environment.
func LoadGlobalFromEnv() (*GlobalConfiguration, error) {
	config := new(GlobalConfiguration)
	if err := loadGlobal(config); err != nil {
		return nil, err
	}
	return config, nil
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Overload(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

func LoadGlobal(filename string) (*GlobalConfiguration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(GlobalConfiguration)
	if err := loadGlobal(config); err != nil {
		return nil, err
	}
	return config, nil
}

func loadGlobal(config *GlobalConfiguration) error {
	// although the package is called "auth" it used to be called "gomess"
	// so environment configs will remain to be called "GOMESS"
	if err := envconfig.Process("gomess", config); err != nil {
		return err
	}

	if err := config.ApplyDefaults(); err != nil {
		return err
	}

	if err := config.Validate(); err != nil {
		return err
	}

	return nil
}

// LoadFile calls godotenv.Load() when the given filename is empty ignoring any
// errors loading, otherwise it calls godotenv.Overload(filename).
//
// godotenv.Load: preserves env, ".env" path is optional
// godotenv.Overload: overrides env, "filename" path must exist
func LoadFile(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Overload(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

// LoadDirectory does nothing when configDir is empty, otherwise it will attempt
// to load a list of configuration files located in configDir by using ReadDir
// to obtain a sorted list of files containing a .env suffix.
//
// When the list is empty it will do nothing, otherwise it passes the file list
// to godotenv.Overload to pull them into the current environment.
func LoadDirectory(configDir string) error {
	if configDir == "" {
		return nil
	}

	// Returns entries sorted by filename
	ents, err := os.ReadDir(configDir)
	if err != nil {
		// We mimic the behavior of LoadGlobal here, if an explicit path is
		// provided we return an error.
		return err
	}

	var paths []string
	for _, ent := range ents {
		if ent.IsDir() {
			continue // ignore directories
		}

		// We only read files ending in .env
		name := ent.Name()
		if !strings.HasSuffix(name, ".env") {
			continue
		}

		// ent.Name() does not include the watch dir.
		paths = append(paths, filepath.Join(configDir, name))
	}

	// If at least one path was found we load the configuration files in the
	// directory. We don't call override without config files because it will
	// override the env vars previously set with a ".env", if one exists.
	if len(paths) > 0 {
		if err := godotenv.Overload(paths...); err != nil {
			return err
		}
	}
	return nil
}
