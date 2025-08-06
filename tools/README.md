# DataCrunch SDK Development Tools

This directory contains development tools and code generation utilities for the DataCrunch SDK.

## 🔧 Service Generator

The service generator creates consistent service client implementations following the established SDK patterns.

### Usage

```bash
# Generate all services
go run tools/cmd/generate/main.go

# Generate specific service only
go run tools/cmd/generate/main.go -service instance

# Dry run (see what would be generated)
go run tools/cmd/generate/main.go -dry-run

# Custom output directory
go run tools/cmd/generate/main.go -output custom/path
```

### Using Makefile

```bash
# Generate all services
make generate-services

# Generate specific service
make generate-service SERVICE=instance

# Dry run
make generate-services-dry-run
```

### Using go generate

```bash
# Run all go:generate directives
go generate ./...
```

## 📁 Directory Structure

```
tools/
├── README.md                 # This file
├── cmd/
│   └── generate/
│       └── main.go          # Service generator command
└── templates/
    └── service.go.tmpl      # Service template file
```

## 🎯 What Gets Generated

Each service gets:

- **Clean service structure** with embedded `*client.Client`
- **ConfigProvider interface** for session integration  
- **Proper client creation** with `New()` and `newClient()` separation
- **Custom initialization hooks** (`initClient`, `initRequest`)
- **Full integration** with credential chain and retry systems
- **Consistent API patterns** across all services

## 🔄 Available Services

| Package | Class | Description |
|---------|-------|-------------|
| `instance` | `Instance` | Instance management |
| `instanceavailability` | `InstanceAvailability` | Instance availability |
| `instancetypes` | `InstanceTypes` | Instance type definitions |
| `locations` | `Locations` | Data center locations |
| `sshkeys` | `SSHKeys` | SSH key management |
| `startscripts` | `StartScripts` | Instance startup scripts |
| `volumes` | `Volumes` | Volume management |
| `volumetypes` | `VolumeTypes` | Volume type definitions |

## 🏗️ Template System

The generator uses Go's `text/template` with embedded files:

- **Template file**: `tools/templates/service.go.tmpl`
- **Variables**: `{{.PackageName}}`, `{{.ClassName}}`, `{{.ServiceName}}`
- **Embedded**: Templates are embedded in the binary for portability

## 🔧 Extending the Generator

### Adding a New Service

1. Add to `services` slice in `tools/cmd/generate/main.go`:
   ```go
   {"newservice", "NewService", "New Service"},
   ```

2. Run the generator:
   ```bash
   go run tools/cmd/generate/main.go -service newservice
   ```

### Modifying the Template  

1. Edit `tools/templates/service.go.tmpl`
2. Regenerate services:
   ```bash
   make generate-services
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