variable "endpoint" {
  description = "Peekaping API endpoint"
  type        = string
  default     = "https://api.peekaping.com"
}

variable "api_key" {
  description = "Peekaping API key for authentication"
  type        = string
  sensitive   = true
}
