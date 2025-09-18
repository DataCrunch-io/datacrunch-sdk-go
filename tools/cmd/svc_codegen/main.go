// Code generation tool for DataCrunch SDK services
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

// Template file definitions
type TemplateFile struct {
	Name     string // Template file name
	Output   string // Output file name
	Required bool   // Whether this file is required for all services
}

// Available template files
var templateFiles = []TemplateFile{
	{"service.go.tmpl", "service.go", true},
	{"api.go.tmpl", "api.go", false},
	{"interface.go.tmpl", "iface/{packagename}.go", false},
	{"integration_test.go.tmpl", "integration_test.go", false},
}

// ServiceConfig holds configuration for each service
type ServiceConfig struct {
	PackageName string // e.g., "instance", "volumes", "sshkeys"
	ClassName   string // e.g., "Instance", "Volumes", "SSHKeys"
	ServiceName string // e.g., "Instance", "Volume", "SSH Key"
}

// Validate checks if the service configuration is valid
func (sc ServiceConfig) Validate() error {
	if sc.PackageName == "" {
		return fmt.Errorf("PackageName cannot be empty")
	}
	if sc.ClassName == "" {
		return fmt.Errorf("ClassName cannot be empty")
	}
	if sc.ServiceName == "" {
		return fmt.Errorf("ServiceName cannot be empty")
	}

	// Package name should be lowercase and contain only letters and numbers
	var alphanumeric = regexp.MustCompile("^[a-z0-9]*$")
	if !alphanumeric.MatchString(sc.PackageName) {
		return fmt.Errorf("PackageName '%s' should contain only lowercase letters and numbers", sc.PackageName)
	}

	// Class name should start with uppercase letter
	if len(sc.ClassName) == 0 || unicode.IsLower(rune(sc.ClassName[0])) {
		return fmt.Errorf("ClassName '%s' should start with an uppercase letter", sc.ClassName)
	}

	return nil
}

// createServiceConfig creates a ServiceConfig from CLI parameters
func createServiceConfig(serviceName, className, displayName string) ServiceConfig {
	// Default className to capitalized serviceName if not provided
	if className == "" {
		if len(serviceName) > 0 {
			className = strings.ToUpper(serviceName[:1]) + serviceName[1:]
		}
	}

	// Default displayName to className if not provided
	if displayName == "" {
		displayName = className
	}

	return ServiceConfig{
		PackageName: serviceName,
		ClassName:   className,
		ServiceName: displayName,
	}
}

func main() {
	var (
		outputDir     = flag.String("output", "service", "Output directory for generated services")
		serviceName   = flag.String("service", "", "Service name to generate (required)")
		className     = flag.String("class", "", "Service class name (optional, defaults to capitalized service name)")
		displayName   = flag.String("name", "", "Service display name (optional, defaults to class name)")
		dryRun        = flag.Bool("dry-run", false, "Show what would be generated without writing files")
		templateDir   = flag.String("templates", "tools/templates", "Directory containing template files")
		generateAll   = flag.Bool("all", true, "Generate all files (service, api, interface, tests)")
		generateAPI   = flag.Bool("api", false, "Generate API files")
		generateIface = flag.Bool("interface", false, "Generate interface files")
		generateTests = flag.Bool("tests", false, "Generate integration test files")
	)
	flag.Parse()

	fmt.Println("🔧 DataCrunch SDK Service Generator")
	fmt.Println("===================================")
	fmt.Println()

	// Validate required parameters
	if *serviceName == "" {
		fmt.Printf("❌ Service name is required. Use -service flag.\n")
		fmt.Printf("Example: go run tools/cmd/svc_codegen/main.go -service myservice\n")
		os.Exit(1)
	}

	// Create service configuration from CLI parameters
	serviceConfig := createServiceConfig(*serviceName, *className, *displayName)

	// Validate service configuration
	if err := serviceConfig.Validate(); err != nil {
		fmt.Printf("❌ Invalid service configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔍 Generating service: %s (%s)\n", serviceConfig.PackageName, serviceConfig.ClassName)
	fmt.Println()

	// Determine which files to generate
	filesToGenerate := []TemplateFile{
		templateFiles[0], // Always generate service.go
	}

	if *generateAll {
		filesToGenerate = templateFiles
	} else {
		if *generateAPI {
			filesToGenerate = append(filesToGenerate, templateFiles[1])
		}
		if *generateIface {
			filesToGenerate = append(filesToGenerate, templateFiles[2])
		}
		if *generateTests {
			filesToGenerate = append(filesToGenerate, templateFiles[3])
		}
	}

	// Load templates
	templates, err := loadTemplates(*templateDir)
	if err != nil {
		fmt.Printf("❌ Failed to load templates: %v\n", err)
		os.Exit(1)
	}

	// Validate templates
	if err := validateTemplates(templates, filesToGenerate); err != nil {
		fmt.Printf("❌ Template validation failed: %v\n", err)
		os.Exit(1)
	}

	// Generate the service
	if err := generateServiceFiles(templates, serviceConfig, *outputDir, *dryRun, filesToGenerate); err != nil {
		fmt.Printf("❌ Failed to generate %s: %v\n", serviceConfig.PackageName, err)
		os.Exit(1)
	}

	fmt.Println()
	if *dryRun {
		fmt.Printf("🏃 Dry run complete! Would generate service: %s\n", serviceConfig.PackageName)
	} else {
		fmt.Printf("✅ Successfully generated service: %s!\n", serviceConfig.PackageName)
	}

	fmt.Println()
	fmt.Println("📦 Generated files include:")
	fmt.Println("  • Clean, simple service structure")
	fmt.Println("  • ConfigProvider interface for session integration")
	fmt.Println("  • Proper client creation with separation of concerns")
	fmt.Println("  • Custom initialization hooks for extensibility")
	fmt.Println("  • Full integration with credential chain and retry system")
}

// loadTemplates loads all template files from the specified directory
func loadTemplates(templateDir string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	for _, tf := range templateFiles {
		templatePath := filepath.Join(templateDir, tf.Name)

		// Check if template file exists
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			if tf.Required {
				return nil, fmt.Errorf("required template file not found: %s", templatePath)
			}
			continue // Skip optional templates that don't exist
		}

		// Load and parse template
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", templatePath, err)
		}

		templates[tf.Name] = tmpl
	}

	return templates, nil
}

// validateTemplates ensures all required templates are loaded and can be executed
func validateTemplates(templates map[string]*template.Template, filesToGenerate []TemplateFile) error {
	// Check that all required templates are loaded
	for _, tf := range filesToGenerate {
		if tf.Required {
			if _, exists := templates[tf.Name]; !exists {
				return fmt.Errorf("required template %s is not loaded", tf.Name)
			}
		}
	}

	// Test template execution with a dummy service config
	testService := ServiceConfig{
		PackageName: "testservice",
		ClassName:   "TestService",
		ServiceName: "Test Service",
	}

	for _, tf := range filesToGenerate {
		tmpl, exists := templates[tf.Name]
		if !exists {
			continue // Skip optional templates that weren't loaded
		}

		// Test template execution without writing to file (discard output)
		var discardBuffer strings.Builder
		if err := tmpl.Execute(&discardBuffer, testService); err != nil {
			return fmt.Errorf("template %s failed validation: %w", tf.Name, err)
		}
	}

	fmt.Printf("✅ All %d templates validated successfully\n", len(templates))
	return nil
}

// generateServiceFiles generates all specified files for a service
func generateServiceFiles(templates map[string]*template.Template, service ServiceConfig, outputDir string, dryRun bool, filesToGenerate []TemplateFile) error {
	serviceDir := filepath.Join(outputDir, service.PackageName)

	if dryRun {
		fmt.Printf("🔍 Would generate service: %s (%s)\n", service.PackageName, service.ClassName)
		for _, tf := range filesToGenerate {
			outputPath := resolveOutputPath(tf.Output, service, serviceDir)
			fmt.Printf("  📄 %s -> %s\n", tf.Name, outputPath)
		}
		return nil
	}

	// Create service directory if it doesn't exist
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", serviceDir, err)
	}

	// Generate each file
	for _, tf := range filesToGenerate {
		tmpl, exists := templates[tf.Name]
		if !exists {
			if tf.Required {
				return fmt.Errorf("required template %s not loaded", tf.Name)
			}
			continue // Skip optional templates that weren't loaded
		}

		if err := generateFile(tmpl, service, serviceDir, tf); err != nil {
			return fmt.Errorf("failed to generate %s: %w", tf.Name, err)
		}
	}

	fmt.Printf("✅ Generated service: %s (%s)\n", service.PackageName, service.ClassName)
	return nil
}

// resolveOutputPath resolves template placeholders in output paths
func resolveOutputPath(outputTemplate string, service ServiceConfig, serviceDir string) string {
	output := strings.ReplaceAll(outputTemplate, "{packagename}", service.PackageName)
	return filepath.Join(serviceDir, output)
}

// generateFile generates a single file from a template
func generateFile(tmpl *template.Template, service ServiceConfig, serviceDir string, tf TemplateFile) error {
	outputPath := resolveOutputPath(tf.Output, service, serviceDir)

	// Create directory for the output file if needed
	outputFileDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputFileDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputFileDir, err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close file %s: %v\n", outputPath, err)
		}
	}()

	// Execute template
	if err := tmpl.Execute(file, service); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("  📄 Generated: %s\n", outputPath)
	return nil
}
