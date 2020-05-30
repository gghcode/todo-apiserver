variable "base_path" {
  type = string
}

variable "bcrypt_cost" {
  type = number
}

variable "graceful_shutdown_timeout_sec" {
  type = number
}

variable "jwt_secret_key" {
  type = string
}

variable "jwt_access_expires_in_sec" {
  type = string
}

variable "jwt_refresh_expires_in_sec" {
  type = string
}

variable "cors_allow_origins" {
  type = string
}

variable "cors_allow_methods" {
  type = string
}

variable "postgres_host" {
  type = string
}

variable "postgres_user" {
  type = string
}

variable "postgres_name" {
  type = string
}

variable "postgres_password" {
  type = string
}

variable "redis_addr" {
  type = string
}
