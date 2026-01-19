package cli

import (
	"bufio"
	"fmt"
	"strings"
)

// TemplateOptions holds the answers to template-specific questions
type TemplateOptions struct {
	// webapp options
	Frontend       string // React, Vue, Svelte, Vanilla
	Backend        string // Go, Node, Python, None
	Authentication bool
	Features       []string

	// cli options
	Language    string // Go, Python, Rust, Node
	Subcommands []string
	ConfigFormat string // YAML, JSON, TOML, None

	// api options
	Database     string // PostgreSQL, MySQL, SQLite, MongoDB, None
	AuthType     string // JWT, OAuth, API Key, None
	Endpoints    []string
	OpenAPISpec  bool

	// library options
	Exports     []string
	CLIWrapper  bool
}

// askTemplateQuestions asks template-specific questions and returns the options
func askTemplateQuestions(reader *bufio.Reader, template string) (*TemplateOptions, error) {
	opts := &TemplateOptions{}

	switch template {
	case "webapp":
		return askWebappQuestions(reader)
	case "cli":
		return askCLIQuestions(reader)
	case "api":
		return askAPIQuestions(reader)
	case "library":
		return askLibraryQuestions(reader)
	default:
		return opts, nil
	}
}

func askWebappQuestions(reader *bufio.Reader) (*TemplateOptions, error) {
	opts := &TemplateOptions{}

	fmt.Println("\n--- Web Application Configuration ---")

	// Frontend framework
	fmt.Println("\nFrontend framework:")
	fmt.Println("  1. React")
	fmt.Println("  2. Vue")
	fmt.Println("  3. Svelte")
	fmt.Println("  4. Vanilla JS")
	fmt.Print("Select [1-4] (default: 1): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	frontends := map[string]string{"1": "React", "2": "Vue", "3": "Svelte", "4": "Vanilla"}
	opts.Frontend = frontends[choice]
	if opts.Frontend == "" {
		opts.Frontend = "React"
	}

	// Backend language
	fmt.Println("\nBackend language:")
	fmt.Println("  1. Go")
	fmt.Println("  2. Node.js")
	fmt.Println("  3. Python")
	fmt.Println("  4. None (frontend only)")
	fmt.Print("Select [1-4] (default: 1): ")
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	backends := map[string]string{"1": "Go", "2": "Node", "3": "Python", "4": "None"}
	opts.Backend = backends[choice]
	if opts.Backend == "" {
		opts.Backend = "Go"
	}

	// Authentication
	fmt.Print("\nInclude authentication? [Y/n]: ")
	auth, _ := reader.ReadString('\n')
	auth = strings.TrimSpace(strings.ToLower(auth))
	opts.Authentication = auth != "n" && auth != "no"

	// Key features
	fmt.Print("\nKey features (comma-separated, e.g., 'dashboard,user-profile,search'): ")
	features, _ := reader.ReadString('\n')
	features = strings.TrimSpace(features)
	if features != "" {
		for _, f := range strings.Split(features, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				opts.Features = append(opts.Features, f)
			}
		}
	}

	return opts, nil
}

func askCLIQuestions(reader *bufio.Reader) (*TemplateOptions, error) {
	opts := &TemplateOptions{}

	fmt.Println("\n--- CLI Tool Configuration ---")

	// Language
	fmt.Println("\nLanguage:")
	fmt.Println("  1. Go")
	fmt.Println("  2. Python")
	fmt.Println("  3. Rust")
	fmt.Println("  4. Node.js")
	fmt.Print("Select [1-4] (default: 1): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	languages := map[string]string{"1": "Go", "2": "Python", "3": "Rust", "4": "Node"}
	opts.Language = languages[choice]
	if opts.Language == "" {
		opts.Language = "Go"
	}

	// Subcommands
	fmt.Print("\nSubcommands (comma-separated, e.g., 'init,run,build'): ")
	subs, _ := reader.ReadString('\n')
	subs = strings.TrimSpace(subs)
	if subs != "" {
		for _, s := range strings.Split(subs, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				opts.Subcommands = append(opts.Subcommands, s)
			}
		}
	}

	// Config format
	fmt.Println("\nConfiguration format:")
	fmt.Println("  1. YAML")
	fmt.Println("  2. JSON")
	fmt.Println("  3. TOML")
	fmt.Println("  4. None")
	fmt.Print("Select [1-4] (default: 1): ")
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	formats := map[string]string{"1": "YAML", "2": "JSON", "3": "TOML", "4": "None"}
	opts.ConfigFormat = formats[choice]
	if opts.ConfigFormat == "" {
		opts.ConfigFormat = "YAML"
	}

	return opts, nil
}

func askAPIQuestions(reader *bufio.Reader) (*TemplateOptions, error) {
	opts := &TemplateOptions{}

	fmt.Println("\n--- REST API Configuration ---")

	// Language
	fmt.Println("\nLanguage:")
	fmt.Println("  1. Go")
	fmt.Println("  2. Python")
	fmt.Println("  3. Node.js")
	fmt.Println("  4. Rust")
	fmt.Print("Select [1-4] (default: 1): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	languages := map[string]string{"1": "Go", "2": "Python", "3": "Node", "4": "Rust"}
	opts.Language = languages[choice]
	if opts.Language == "" {
		opts.Language = "Go"
	}

	// Database
	fmt.Println("\nDatabase:")
	fmt.Println("  1. PostgreSQL")
	fmt.Println("  2. MySQL")
	fmt.Println("  3. SQLite")
	fmt.Println("  4. MongoDB")
	fmt.Println("  5. None")
	fmt.Print("Select [1-5] (default: 1): ")
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	databases := map[string]string{"1": "PostgreSQL", "2": "MySQL", "3": "SQLite", "4": "MongoDB", "5": "None"}
	opts.Database = databases[choice]
	if opts.Database == "" {
		opts.Database = "PostgreSQL"
	}

	// Auth type
	fmt.Println("\nAuthentication type:")
	fmt.Println("  1. JWT")
	fmt.Println("  2. OAuth")
	fmt.Println("  3. API Key")
	fmt.Println("  4. None")
	fmt.Print("Select [1-4] (default: 1): ")
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	authTypes := map[string]string{"1": "JWT", "2": "OAuth", "3": "API Key", "4": "None"}
	opts.AuthType = authTypes[choice]
	if opts.AuthType == "" {
		opts.AuthType = "JWT"
	}

	// Key endpoints
	fmt.Print("\nKey endpoints (comma-separated, e.g., 'users,products,orders'): ")
	endpoints, _ := reader.ReadString('\n')
	endpoints = strings.TrimSpace(endpoints)
	if endpoints != "" {
		for _, e := range strings.Split(endpoints, ",") {
			e = strings.TrimSpace(e)
			if e != "" {
				opts.Endpoints = append(opts.Endpoints, e)
			}
		}
	}

	// OpenAPI spec
	fmt.Print("\nGenerate OpenAPI spec? [Y/n]: ")
	spec, _ := reader.ReadString('\n')
	spec = strings.TrimSpace(strings.ToLower(spec))
	opts.OpenAPISpec = spec != "n" && spec != "no"

	return opts, nil
}

func askLibraryQuestions(reader *bufio.Reader) (*TemplateOptions, error) {
	opts := &TemplateOptions{}

	fmt.Println("\n--- Library Configuration ---")

	// Language
	fmt.Println("\nLanguage:")
	fmt.Println("  1. Go")
	fmt.Println("  2. Python")
	fmt.Println("  3. Node.js (npm)")
	fmt.Println("  4. Rust")
	fmt.Print("Select [1-4] (default: 1): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	languages := map[string]string{"1": "Go", "2": "Python", "3": "Node", "4": "Rust"}
	opts.Language = languages[choice]
	if opts.Language == "" {
		opts.Language = "Go"
	}

	// Key exports
	fmt.Print("\nKey exports/functions (comma-separated, e.g., 'Parse,Format,Validate'): ")
	exports, _ := reader.ReadString('\n')
	exports = strings.TrimSpace(exports)
	if exports != "" {
		for _, e := range strings.Split(exports, ",") {
			e = strings.TrimSpace(e)
			if e != "" {
				opts.Exports = append(opts.Exports, e)
			}
		}
	}

	// CLI wrapper
	fmt.Print("\nInclude CLI wrapper? [y/N]: ")
	cli, _ := reader.ReadString('\n')
	cli = strings.TrimSpace(strings.ToLower(cli))
	opts.CLIWrapper = cli == "y" || cli == "yes"

	return opts, nil
}

// ToMap converts TemplateOptions to a map for JSON serialization
func (o *TemplateOptions) ToMap(template string) map[string]interface{} {
	result := make(map[string]interface{})

	switch template {
	case "webapp":
		result["frontend"] = o.Frontend
		result["backend"] = o.Backend
		result["authentication"] = o.Authentication
		result["features"] = o.Features
	case "cli":
		result["language"] = o.Language
		result["subcommands"] = o.Subcommands
		result["config_format"] = o.ConfigFormat
	case "api":
		result["language"] = o.Language
		result["database"] = o.Database
		result["auth_type"] = o.AuthType
		result["endpoints"] = o.Endpoints
		result["openapi_spec"] = o.OpenAPISpec
	case "library":
		result["language"] = o.Language
		result["exports"] = o.Exports
		result["cli_wrapper"] = o.CLIWrapper
	}

	return result
}
