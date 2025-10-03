# Gym Local Server

A Go-based local server for managing IP cameras with PTZ (Pan-Tilt-Zoom) control and WebRTC streaming capabilities. This server provides a RESTful API for camera management and real-time video streaming in gym environments.

## Features

- **Multiple Camera Support**: Manage multiple IP cameras simultaneously
- **PTZ Control**: Control camera pan, tilt, and zoom operations via ONVIF protocol
- **WebRTC Streaming**: Real-time video streaming using WebRTC technology
- **ONVIF Compatible**: Works with ONVIF-compliant IP cameras
- **RESTful API**: Clean and documented REST API endpoints

## Requirements

- Go 1.24.1 or higher
- ONVIF-compatible IP cameras
- Network access to cameras

## Installation

### Clone the repository

```bash
git clone https://github.com/Illia-33/gym-localserver.git
cd gym-localserver
```

### Install dependencies

```bash
go mod download
```

### Build the project

```bash
make all
```

This will build two executables:
- `configurator`: Tool for generating configuration files
- `server`: The main local server application

## Configuration

### Generate Configuration

Create a `config.yml` file using the configurator:

```bash
make config
```

Or manually create a configuration file with the following structure:

```yaml
settings:
  auth_key: "your-authentication-key"

cameras:
  - label: "Camera 1"
    description: "Front entrance camera"
    type: "onvif"
    ip: "192.168.1.100"
    port: 80
    login: "admin"
    password: "password123"
  
  - label: "Camera 2"
    description: "Main gym area"
    type: "onvif"
    ip: "192.168.1.101"
    port: 80
    login: "admin"
    password: "password123"
```

### Configuration Parameters

- **settings.auth_key**: Authentication key for gym central server registration
- **cameras**: Array of camera configurations
  - **label**: Unique identifier for the camera
  - **description**: Human-readable description
  - **type**: Camera type (currently supports "onvif")
  - **ip**: Camera IP address
  - **port**: Camera port (typically 80 or 8080)
  - **login**: Camera username
  - **password**: Camera password

## Usage

### Start the Server

Run the server with default settings:

```bash
./server
```

Or specify custom bind address and config file:

```bash
./server -bind 0.0.0.0:8080 -config /path/to/config.yml
```

### Command Line Options

- `-bind`: Address to bind the server on (default: `0.0.0.0:8080`)
- `-config`: Path to configuration file (default: `./config.yml`)

### Quick Start

Build, configure, and run everything in one command:

```bash
make run
```

## API Documentation

The server exposes a RESTful API at `/api/v1`. Full OpenAPI/Swagger documentation is available in [`docs/client.yml`](docs/client.yml).

### Base URL

```
http://localhost:8080/api/v1
```

### Endpoints

#### Get Cameras List

Get information about all available cameras.

```http
GET /api/v1/cameras
```

**Response:**
```json
{
  "cameras": [
    {
      "label": "Camera 1",
      "description": "Front entrance camera"
    },
    {
      "label": "Camera 2",
      "description": "Main gym area"
    }
  ]
}
```

#### Start PTZ Control

Start PTZ (Pan-Tilt-Zoom) control for a specific camera.

```http
POST /api/v1/camera/{camera_id}/ptz
```

**Parameters:**
- `camera_id` (path): Camera index (0-based)

**Request Body:**
```json
{
  "velocity": {
    "pan": 0.5,
    "tilt": 0.3,
    "zoom": 0.0
  },
  "deadline": "5s"
}
```

**Velocity ranges:**
- `pan`: -1.0 to 1.0 (left to right)
- `tilt`: -1.0 to 1.0 (down to up)
- `zoom`: -1.0 to 1.0 (zoom out to zoom in)
- `deadline`: Duration string (e.g., "5s", "1m")

#### Stop PTZ Control

Stop PTZ control for a specific camera.

```http
DELETE /api/v1/camera/{camera_id}/ptz
```

**Parameters:**
- `camera_id` (path): Camera index (0-based)

#### Setup WebRTC Connection

Establish a WebRTC connection for video streaming.

```http
POST /api/v1/camera/{camera_id}/webrtc
```

**Parameters:**
- `camera_id` (path): Camera index (0-based)

**Request Body:**
```json
{
  "offer_b64": "base64-encoded-sdp-offer"
}
```

**Response:**
```json
{
  "id": "connection-id",
  "local_desc_b64": "base64-encoded-sdp-answer"
}
```

## Project Structure

```
.
├── api/                    # API definitions
│   └── localserver/       # Local server API types
├── cmd/                   # Command-line applications
│   ├── configurator/      # Configuration tool
│   └── localserver/       # Main server application
├── docs/                  # Documentation
│   └── client.yml         # OpenAPI/Swagger specification
├── internal/              # Internal packages
│   ├── config/           # Configuration handling
│   ├── localserver/      # Server implementation
│   │   ├── server/       # HTTP server and API handlers
│   │   └── service/      # Business logic
│   └── webrtc/           # WebRTC implementation
├── pkg/                   # Public packages
│   ├── camera/           # Camera abstractions
│   ├── config/           # Configuration types
│   ├── onvif/            # ONVIF protocol implementation
│   ├── rtsp/             # RTSP protocol client implementation
│   └── sdp/              # SDP parsing
├── Makefile              # Build automation
├── go.mod                # Go module definition
└── README.md             # This file
```

## Development

### Build Commands

```bash
# Build all binaries
make all

# Build configurator only
make configurator

# Build server only
make server

# Generate configuration
make config

# Build and run
make run
```

### Running Tests

```bash
go test ./...
```

## Dependencies

Key dependencies include:

- **Gin**: HTTP web framework
- **Pion WebRTC**: WebRTC implementation for Go
- **ONVIF**: ONVIF protocol support
- **YAML**: Configuration file parsing

See [`go.mod`](go.mod) for the complete list of dependencies.

## Architecture

The server follows a clean architecture pattern:

1. **API Layer** (`internal/localserver/server`): HTTP handlers and routing
2. **Service Layer** (`internal/localserver/service`): Business logic and orchestration
3. **Domain Layer** (`pkg/camera`): Core camera abstractions
4. **Infrastructure Layer** (`pkg/onvif`, `pkg/rtsp`): External protocol implementations

## Troubleshooting

### Camera Not Responding

- Verify camera IP address and port
- Ensure camera is ONVIF-compliant
- Check network connectivity
- Verify login credentials

### WebRTC Connection Issues

- Ensure proper NAT/firewall configuration
- Check that camera supports RTSP streaming
- Verify SDP offer format is correct

### Server Won't Start

- Check if port 8080 is already in use
- Verify config.yml exists and is valid
- Ensure at least one camera is accessible

## Acknowledgments

- Built with [Gin Web Framework](https://github.com/gin-gonic/gin)
- WebRTC powered by [Pion](https://github.com/pion/webrtc)
- ONVIF support via [use-go/onvif](https://github.com/use-go/onvif)
