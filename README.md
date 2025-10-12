<div align="center">

# 🚀 Terraform Provider for Peekaping

**The most comprehensive Terraform provider for Peekaping monitoring infrastructure**

[![Terraform Registry](https://img.shields.io/badge/Terraform%20Registry-7B42BC?style=for-the-badge&logo=terraform&logoColor=white)](https://registry.terraform.io/providers/tafaust/peekaping/latest)
[![GitHub stars](https://img.shields.io/github/stars/tafaust/terraform-provider-peekaping?style=for-the-badge&logo=github)](https://github.com/tafaust/terraform-provider-peekaping/stargazers)
[![Support on Ko-fi](https://img.shields.io/badge/Support%20on-Ko--fi-FF5E5B?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/tafaust)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-0.1.0-blue?style=for-the-badge)](https://github.com/tafaust/terraform-provider-peekaping/releases)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)

[![Monitor Types](https://img.shields.io/badge/Monitor%20Types-18-brightgreen?style=for-the-badge)](#monitors)
[![Resources](https://img.shields.io/badge/Resources-6-purple?style=for-the-badge)](#resources-supported)
[![Peekaping](https://img.shields.io/badge/Peekaping-0.0.39-orange?style=for-the-badge)](https://peekaping.com)

</div>

## 🌟 Show Your Support

<div align="center">

**If this project helped you, please consider giving it a ⭐!**

[![GitHub stars](https://img.shields.io/github/stars/tafaust/terraform-provider-peekaping?style=social&label=Star)](https://github.com/tafaust/terraform-provider-peekaping/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tafaust/terraform-provider-peekaping?style=social&label=Fork)](https://github.com/tafaust/terraform-provider-peekaping/network)
[![GitHub watchers](https://img.shields.io/github/watchers/tafaust/terraform-provider-peekaping?style=social&label=Watch)](https://github.com/tafaust/terraform-provider-peekaping/watchers)

</div>

---

## ✨ Features & Capabilities

### 🚀 Key Features

- **🎯 Complete Infrastructure as Code** - Full CRUD operations, state management, and import support
- **🔐 Enterprise Security** - 2FA support, secure credentials, and environment variable management
- **📊 Rich Monitoring** - 18 monitor types with JSON configuration and advanced scheduling
- **🎨 Status Pages** - Public status pages with custom themes and real-time updates
- **🔔 Smart Notifications** - Multiple channels with flexible routing and default settings
- **🏷️ Organization** - Tag-based organization, maintenance windows, and proxy support

### 🔧 Resources

- **peekaping_monitor** - Create and manage monitoring checks (HTTP, HTTP-keyword, HTTP-JSON-query, TCP, Ping, DNS, Push, Docker, gRPC-keyword, SNMP, MongoDB, MySQL, PostgreSQL, SQL Server, Redis, MQTT, RabbitMQ, Kafka Producer)
- **peekaping_notification** - Configure notification channels (Email, webhook, custom channels)
- **peekaping_tag** - Organize monitors with tags (All monitor types)
- **peekaping_maintenance** - Schedule maintenance windows (All monitor types)
- **peekaping_status_page** - Create public status pages (All monitor types)
- **peekaping_proxy** - Configure proxy settings (All monitor types)

### 📋 Data Sources

- **peekaping_monitor** - Query existing monitors (Find monitors by name, type, or tags)
- **peekaping_notification** - Query notification channels (Discover available notification channels)
- **peekaping_tag** - Query tags (Lookup tags and their configurations)
- **peekaping_maintenance** - Query maintenance windows (Find scheduled maintenance periods)
- **peekaping_status_page** - Query status pages (Discover existing status pages)
- **peekaping_proxy** - Query proxy configurations (Find available proxy settings)

## 📦 Installation & Setup

### Option 1: Terraform Registry (Recommended)

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.1.0"
    }
  }
}
```

Then run:
```bash
terraform init
```

### Option 2: Local Development

For development or testing with the latest changes:

1. **Clone and build**:
   ```bash
   git clone https://github.com/tafaust/terraform-provider-peekaping
   cd terraform-provider-peekaping
   go build -o terraform-provider-peekaping
   ```

2. **Configure Terraform** (`~/.terraformrc`):
   ```hcl
   provider_installation {
     dev_overrides {
       "tafaust/peekaping" = "/path/to/terraform-provider-peekaping"
     }
     direct {}
   }
   ```

3. **Initialize**:
   ```bash
   terraform init
   ```

### 🔧 Provider Configuration

Configure the provider in your Terraform files:

```hcl
provider "peekaping" {
  endpoint = "https://api.peekaping.com"  # Your Peekaping instance URL
  email    = "your-email@example.com"     # Your Peekaping email
  password = "your-password"              # Your Peekaping password
  token    = "123456"                     # 2FA token (if 2FA is enabled)
}
```

**Environment Variables** (Alternative to provider config):
```bash
export PEEKAPING_ENDPOINT="https://api.peekaping.com"
export PEEKAPING_EMAIL="your-email@example.com"
export PEEKAPING_PASSWORD="your-password"
export PEEKAPING_TOKEN="123456"  # 2FA token
```

## 📚 Examples & Use Cases

### 🎯 Real-World Examples

| Example | Description | Complexity |
|---------|-------------|------------|
| **[Simple](examples/simple/)** | Basic HTTP monitor setup | ⭐ Beginner |
| **[Comprehensive](examples/comprehensive/)** | Full monitoring stack | ⭐⭐⭐ Advanced |
| **[Full Lifecycle](examples/full-lifecycle/)** | Complete workflow demo | ⭐⭐ Intermediate |

### 🚀 Quick Examples

**HTTP Monitor with Notifications:**
```hcl
resource "peekaping_monitor" "api" {
  name = "API Health Check"
  type = "http"
  config = jsonencode({
    url = "https://api.example.com/health"
    method = "GET"
    headers = {
      "Authorization" = "Bearer ${var.api_token}"
    }
  })
  interval = 30
  timeout  = 10
}

resource "peekaping_notification" "email" {
  name = "Team Alerts"
  type = "email"
  config = jsonencode({
    emails = ["team@example.com"]
  })
}
```

**Database Monitoring:**
```hcl
resource "peekaping_monitor" "postgres" {
  name = "PostgreSQL Database"
  type = "postgresql"
  config = jsonencode({
    hostname = "db.example.com"
    port     = 5432
    username = "monitor"
    password = var.db_password
    database = "app_db"
  })
}
```

## 🧪 Testing & Validation

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/provider -run TestResourceMonitor
```

### Example Validation

```bash
# Validate examples
cd examples/simple
terraform init
terraform plan
terraform apply

# Test comprehensive setup
cd examples/comprehensive
terraform init
terraform plan
terraform apply
```

## 🛠️ Development

### Project Structure

```
├── internal/
│   ├── peekaping/          # API client implementation
│   └── provider/           # Terraform provider resources
├── examples/               # Real-world usage examples
├── docs/                   # Comprehensive documentation
├── tools/                  # Development tools
└── main.go                 # Provider entry point
```

### Building

```bash
# Build the provider
go build -o terraform-provider-peekaping

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o terraform-provider-peekaping-linux
GOOS=windows GOARCH=amd64 go build -o terraform-provider-peekaping-windows.exe
```

### Development Workflow

```bash
# 1. Make changes to the code
# 2. Run tests
go test ./...

# 3. Build the provider
go build -o terraform-provider-peekaping

# 4. Test with examples
cd examples/simple
terraform init
terraform plan
terraform apply
```

## 📝 Known Limitations

- **2FA Token Expiration**: 2FA tokens expire quickly and need frequent updates during development
- **JSON Formatting**: Monitor config JSON is normalized for consistency across different monitor types

## 🤝 Contributing

We **love** contributions! This project thrives on community input and collaboration.

### 🌟 Ways to Contribute

| Contribution Type | Description | Impact |
|------------------|-------------|---------|
| 🐛 **Bug Reports** | Found an issue? Let us know! | High |
| 💡 **Feature Requests** | Have an idea? We'd love to hear it! | High |
| 📝 **Documentation** | Improve examples, docs, or README | Medium |
| 🔧 **Code Contributions** | Fix bugs or add features | High |
| 🧪 **Testing** | Add tests or improve coverage | Medium |
| 📢 **Community** | Help others, answer questions | High |

### 🚀 Getting Started

1. **⭐ Star this repository** - It helps others discover the project!
2. **🍴 Fork the repository**
3. **🌿 Create a feature branch**: `git checkout -b feature/amazing-feature`
4. **💻 Make your changes**
5. **✅ Add tests** for new functionality
6. **📝 Update documentation** if needed
7. **🔄 Submit a pull request**

### 🎯 Contribution Guidelines

- **Code Style**: Follow Go conventions and run `gofmt`
- **Testing**: Add tests for new features and bug fixes
- **Documentation**: Update README and examples as needed
- **Commit Messages**: Use clear, descriptive commit messages
- **Pull Requests**: Provide a clear description of changes

### 💬 Community

- **Discussions**: Use GitHub Discussions for questions and ideas
- **Issues**: Report bugs and request features
- **Pull Requests**: Submit code improvements

**Every contribution matters!** Whether you're fixing a typo or adding a new monitor type, your help makes this project better for everyone.

## 🆘 Support & Community

| Resource | Description | Link |
|----------|-------------|------|
| 📚 **Documentation** | Complete API reference and guides | [docs/](docs/) |
| 💡 **Examples** | Real-world usage examples | [examples/](examples/) |
| 🐛 **Bug Reports** | Found an issue? Let us know! | [GitHub Issues](https://github.com/tafaust/terraform-provider-peekaping/issues) |
| 💬 **Discussions** | Questions and community chat | [GitHub Discussions](https://github.com/tafaust/terraform-provider-peekaping/discussions) |
| 📖 **Terraform Registry** | Official provider documentation | [Registry](https://registry.terraform.io/providers/tafaust/peekaping/latest) |

## 🎯 Roadmap & Status

### ✅ Completed Features
- [x] **All Core Resources** - Monitors, notifications, tags, maintenance, status pages, proxies
- [x] **Data Sources** - Query existing resources
- [x] **18 Monitor Types** - Comprehensive monitoring coverage
- [x] **2FA Support** - Enhanced security
- [x] **Terraform Registry** - Official provider distribution
- [x] **Comprehensive Examples** - Real-world usage patterns
- [x] **Production Testing** - Battle-tested reliability

### 🚧 In Progress
- [ ] **API Key Authentication** - Alternative auth method ([PR pending](https://github.com/0xfurai/peekaping/pull/204))
- [ ] **Enhanced Error Handling** - Better user experience

---

<div align="center">

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

**Made with ❤️ by the community, for the community**

[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Powered by Terraform](https://img.shields.io/badge/Powered%20by-Terraform-7B42BC?style=for-the-badge&logo=terraform)](https://terraform.io)

</div>
