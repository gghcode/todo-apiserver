variable "image_name" {
  description = "Name of the docker image to deploy."
  default     = "gcr.io/vapias/todo-apiserver"
}

variable "image_tag" {
  description = "The docker image tag to deploy."
}

variable "service_name" {
  description = "Name of cloud run service"
  default     = "todo-apiserver"
}

variable "service_location" {
  description = "Location where placed service"
  default     = "asia-northeast1"
}

variable "project" {
  description = "Name of Google cloud platform project"
  default     = "vapias"
}
