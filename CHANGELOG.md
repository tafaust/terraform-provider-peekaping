# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2025-11-20

### Fixed
- **State Drift** - Fixed bug where provider ignored `tag_ids` and `notification_ids` from API responses, causing state drift and "inconsistent final plan" errors
- Updated `setModelFromMonitor` and `setModelFromMonitorWithState` to populate `tag_ids` and `notification_ids` from API responses
- Marked `tag_ids` and `notification_ids` as `Computed` in schema to allow Terraform to sync values from API

## [0.2.0] - 2025-11-20

### Fixed
- **API Key Authentication** - Fixed API key header format (removed Bearer prefix requirement)

### Changed
- Updated Peekaping version reference to 0.0.41 in documentation

## [0.1.1] - 2025-10-12

### Added
- **Default values** - Automatic defaults for monitor fields (interval: 60s, timeout: 30s, max_retries: 3, retry_interval: 60s, resend_interval: 10s, active: true)
- **Enhanced validation** - Connection string format examples and validation guidance for database monitors
- **Version management script** - Automated version bumping via `make version-bump NEW_VERSION=x.y.z`

### Changed
- **Documentation** - Consolidated monitor configuration documentation into a single comprehensive guide
- **Field requirements** - `notification_ids` is now properly marked as required

### Fixed
- **Timeout validation** - Added constraint that timeout must be less than 80% of interval
- **Default values** - Source code now properly implements default values matching documentation
- **Documentation accuracy** - All monitor configuration examples now match the actual API requirements
- Added `retry_interval` field requirement to prevent API validation errors
- Fixed formatting inconsistencies in monitor configuration examples

## [0.1.0] - 2025-10-12

### Added
- **API Key Authentication** - Primary authentication method using API keys (recommended)
- Initial release of the Peekaping Terraform Provider
- Support for managing Peekaping monitors (HTTP, TCP, Ping, DNS, Push, gRPC)
- Support for managing notification channels (Email, Webhook)
- Support for managing tags for organizing monitors
- Support for managing proxy configurations
- Support for managing maintenance windows (one-time, recurring, cron)
- Support for managing status pages with custom themes and styling
- Comprehensive data sources for querying existing resources
- Full CRUD operations for all resources
- Input validation and error handling
- Comprehensive documentation and examples

### Changed
- **Authentication method** - API key authentication is now the recommended method
- **Field naming** - Renamed `token` to `totp_token` to clarify it's for 2FA authentication

## [0.0.3] - 2025-10-05

### Changed
- Updated documentation and added new data sources

## [0.0.2] - 2025-10-05

### Changed
- Updated documentation and GoReleaser configuration

## [0.0.1] - 2025-10-05

### Added
- Initial release
- Complete implementation of all Peekaping API resources
- Comprehensive testing suite
- Full documentation
