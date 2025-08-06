// Code generation tool for DataCrunch SDK services
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Templates embedded at build time
var serviceTemplate = `package {{.PackageName}}

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/restjson"
)

const (
	EndpointsID = "{{.PackageName}}"
	APIVersion  = "v1"
)

// {{.ClassName}} provides the API operation methods for making requests to
// DataCrunch {{.ServiceName}} API
type {{.ClassName}} struct {
	*client.Client
}

// Client is an alias for {{.ClassName}} to match the expected interface
type Client = *{{.ClassName}}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

type ConfigProvider interface {
	ClientConfig(serviceName string, cfgs ...*interface{}) client.Config
}

// New creates a new instance of the {{.ClassName}} client with a config provider.
// If additional configuration is needed for the client instance use the optional
// client.Config parameter to add your extra config.
//
// Example:
//
//	mySession := session.Must(session.New())
//
//	// Create a {{.ClassName}} client from just a session.
//	svc := {{.PackageName}}.New(mySession)
//
//	// Create a {{.ClassName}} client with additional configuration
//	svc := {{.PackageName}}.New(mySession, &client.Config{Timeout: 60 * time.Second})
func New(p ConfigProvider, cfgs ...*interface{}) *{{.ClassName}} {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(c)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg client.Config) *{{.ClassName}} {
	handlers := request.Handlers{}

	// Add protocol handlers for REST JSON
	handlers.Build.PushBackNamed(restjson.BuildHandler)
	handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	handlers.Complete.PushBackNamed(restjson.UnmarshalMetaHandler)

	svc := &{{.ClassName}}{
		Client: client.New(&cfg, metadata.ClientInfo{
			ServiceName: EndpointsID,
			APIVersion:  APIVersion,
			Endpoint:    cfg.BaseURL,
		}, handlers),
	}

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

func (c *{{.ClassName}}) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}`

// ServiceConfig holds configuration for each service
type ServiceConfig struct {
	PackageName string // e.g., "instance", "volumes", "sshkeys"
	ClassName   string // e.g., "Instance", "Volumes", "SSHKeys"
	ServiceName string // e.g., "Instance", "Volume", "SSH Key"
}

// All DataCrunch SDK services
var services = []ServiceConfig{
	{"instance", "Instance", "Instance"},
	{"instanceavailability", "InstanceAvailability", "Instance Availability"},
	{"instancetypes", "InstanceTypes", "Instance Types"},
	{"locations", "Locations", "Locations"},
	{"sshkeys", "SSHKeys", "SSH Keys"},
	{"startscripts", "StartScripts", "Start Scripts"},
	{"volumes", "Volumes", "Volumes"},
	{"volumetypes", "VolumeTypes", "Volume Types"},
}

func main() {
	var (
		outputDir = flag.String("output", "service", "Output directory for generated services")
		service   = flag.String("service", "", "Generate only a specific service (optional)")
		dryRun    = flag.Bool("dry-run", false, "Show what would be generated without writing files")
	)
	flag.Parse()

	fmt.Println("🔧 DataCrunch SDK Service Generator")
	fmt.Println("===================================")
	fmt.Println()

	// Load template
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		fmt.Printf("❌ Failed to parse template: %v\n", err)
		os.Exit(1)
	}

	// Filter services if specific service requested
	servicesToGenerate := services
	if *service != "" {
		servicesToGenerate = nil
		for _, svc := range services {
			if svc.PackageName == *service || svc.ClassName == *service {
				servicesToGenerate = []ServiceConfig{svc}
				break
			}
		}
		if len(servicesToGenerate) == 0 {
			fmt.Printf("❌ Service '%s' not found. Available services:\n", *service)
			for _, svc := range services {
				fmt.Printf("  - %s (%s)\n", svc.PackageName, svc.ClassName)
			}
			os.Exit(1)
		}
	}

	// Generate services
	generated := 0
	for _, svc := range servicesToGenerate {
		if err := generateService(tmpl, svc, *outputDir, *dryRun); err != nil {
			fmt.Printf("❌ Failed to generate %s: %v\n", svc.PackageName, err)
			continue
		}
		generated++
	}

	fmt.Println()
	if *dryRun {
		fmt.Printf("🏃 Dry run complete! Would generate %d service(s)\n", generated)
	} else {
		fmt.Printf("✅ Successfully generated %d service(s)!\n", generated)
	}

	if generated > 0 {
		fmt.Println()
		fmt.Println("📦 Generated services include:")
		fmt.Println("  • Clean, simple service structure")
		fmt.Println("  • ConfigProvider interface for session integration")
		fmt.Println("  • Proper client creation with separation of concerns")
		fmt.Println("  • Custom initialization hooks for extensibility")
		fmt.Println("  • Full integration with credential chain and retry system")
	}
}

func generateService(tmpl *template.Template, service ServiceConfig, outputDir string, dryRun bool) error {
	// Create service directory path
	serviceDir := filepath.Join(outputDir, service.PackageName)
	servicePath := filepath.Join(serviceDir, "service.go")

	if dryRun {
		fmt.Printf("🔍 Would generate: %s -> %s (%s)\n", service.PackageName, servicePath, service.ClassName)
		return nil
	}

	// Create service directory if it doesn't exist
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", serviceDir, err)
	}

	// Create service file
	file, err := os.Create(servicePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", servicePath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, service); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("✅ Generated: %s -> %s (%s)\n", service.PackageName, servicePath, service.ClassName)
	return nil
}
