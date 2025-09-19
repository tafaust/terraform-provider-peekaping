variable "endpoint" {
  description = "Peekaping API endpoint"
  type        = string
  default     = "https://api.peekaping.com"
}

variable "email" {
  description = "Peekaping account email"
  type        = string
}

variable "password" {
  description = "Peekaping account password"
  type        = string
  sensitive   = true
}
