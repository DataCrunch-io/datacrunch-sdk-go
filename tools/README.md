# DataCrunch SDK Development Tools

This directory contains development tools and code generation utilities for the DataCrunch SDK.

## 🔧 Service Generator

The service generator creates consistent service client implementations following the established SDK patterns.

### Usage

```bash
# Generate a new service (required: service name)
go run tools/cmd/svc_codegen/main.go -service myservice

# Generate with custom class and display names
go run tools/cmd/svc_codegen/main.go -service myservice -class MyService -name "My Service"

# Generate specific file types only
go run tools/cmd/svc_codegen/main.go -service myservice -api -interface

# Dry run (see what would be generated)
go run tools/cmd/svc_codegen/main.go -service myservice -dry-run

# Custom output directory
go run tools/cmd/svc_codegen/main.go -service myservice -output custom/path
```

### Using Makefile (Recommended)

```bash
# Generate a new service (required: SERVICE parameter)
make generate-service SERVICE=myservice

# Generate with custom class and display names
make generate-service SERVICE=myservice CLASS=MyService NAME="My Service"

# Dry run mode
make generate-service SERVICE=myservice DRY_RUN=true
```

### Using go generate

```bash
# Run all go:generate directives
go generate ./...
```

## 📁 Directory Structure

```
tools/
├── README.md                      # This file
├── cmd/
│   └── svc_codegen/
│       └── main.go               # Service generator command
└── templates/
    ├── service.go.tmpl           # Service client template
    ├── api.go.tmpl               # API methods template
    ├── interface.go.tmpl         # Interface definition template
    └── integration_test.go.tmpl  # Integration test template
```

## 🎯 What Gets Generated

Each service gets a complete set of files:

- **service.go** - Clean service structure with embedded `*client.Client`
- **api.go** - API method definitions with request/response types
- **iface/[service].go** - Interface definitions for the service
- **integration_test.go** - Integration test templates

Features included:
- **ConfigProvider interface** for session integration
- **Proper client creation** with `New()` and `newClient()` separation
- **Custom initialization hooks** (`initClient`, `initRequest`)
- **Full integration** with credential chain and retry systems
- **Consistent API patterns** across all services

## 🔄 Service Generation

The generator creates services dynamically based on the parameters you provide. You can generate any service by specifying:

- **Service Name** (required): The package name (e.g., `myservice`)
- **Class Name** (optional): The Go struct name (defaults to capitalized service name)
- **Display Name** (optional): Human-readable name for documentation (defaults to class name)

## 🏗️ Template System

The generator uses Go's `text/template` with embedded files:

- **Template file**: `tools/templates/service.go.tmpl`
- **Variables**: `{{.PackageName}}`, `{{.ClassName}}`, `{{.ServiceName}}`
- **Embedded**: Templates are embedded in the binary for portability

## 🔧 Extending the Generator

### Generating a New Service

Simply run the generator with your desired service name:

```bash
make generate-service SERVICE=myservice
```

The generator will automatically create all necessary files with proper naming conventions.

### Modifying Templates

1. Edit templates in `tools/templates/` directory
2. Regenerate your service to test changes:

   ```bash
   make generate-service SERVICE=myservice DRY_RUN=true
   ```

### Custom Templates

Create additional `.tmpl` files in `tools/templates/` and modify the generator to use them.

## 🚀 Best Practices

1. **Always use the generator** for new services to maintain consistency
2. **Regenerate after template changes** to keep all services in sync
3. **Use dry-run mode** to preview changes before applying
4. **Version control templates** as they define the SDK's architecture
5. **Test generated code** to ensure it compiles and works correctly

The service generator ensures all DataCrunch SDK services follow the same patterns and get automatic updates when the underlying architecture evolves! 🎯
