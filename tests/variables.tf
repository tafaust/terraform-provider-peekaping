# Variables for Peekaping Monitor testing

variable "monitor_name" {
  description = "Name of the monitor to create"
  type        = string
  default     = "test-monitor"
}

variable "monitor_type" {
  description = "Type of monitor to create"
  type        = string
  default     = "http"
}

variable "monitor_url" {
  description = "URL for HTTP monitors"
  type        = string
  default     = "https://example.com"
}

variable "monitor_host" {
  description = "Host for TCP, Ping, and DNS monitors"
  type        = string
  default     = "example.com"
}

variable "monitor_port" {
  description = "Port for TCP monitors"
  type        = number
  default     = 80
}

variable "packet_size" {
  description = "Packet size for Ping monitors"
  type        = number
  default     = 56
}

variable "resolver_server" {
  description = "DNS resolver server for DNS monitors"
  type        = string
  default     = "8.8.8.8"
}

variable "resolver_port" {
  description = "DNS resolver port for DNS monitors"
  type        = number
  default     = 53
}

variable "resolve_type" {
  description = "DNS record type for DNS monitors"
  type        = string
  default     = "A"
}

variable "connection_string" {
  description = "Database connection string for database monitors"
  type        = string
  default     = "mysql://user:password@localhost:3306/testdb"
}

variable "monitor_interval" {
  description = "Monitor check interval in seconds"
  type        = number
  default     = null
}

variable "monitor_timeout" {
  description = "Monitor timeout in seconds"
  type        = number
  default     = null
}

