provider "google" {
  project = "vapias"
  region  = "asia-northeast1"

  version = "~> v3.19.0"
}

terraform {
  backend "remote" {
    organization = "vapias"
    workspaces {
      name = "todo-apiserver"
    }
  }
}

resource "google_project_service" "run" {
  project                    = var.project
  service                    = "run.googleapis.com"
  disable_dependent_services = false
  disable_on_destroy         = true
}

resource "google_cloud_run_service" "this" {
  name     = var.service_name
  location = var.service_location

  metadata {
    namespace = var.project
  }

  template {
    spec {
      container_concurrency = 0

      containers {
        image = "${var.image_name}:${var.image_tag}"
        env {
          name  = "REST_POSTGRES_HOST"
          value = var.postgres_host
        }
        env {
          name  = "REST_POSTGRES_USER"
          value = var.postgres_user
        }
        env {
          name  = "REST_POSTGRES_NAME"
          value = var.postgres_name
        }
        env {
          name  = "REST_POSTGRES_PASSWORD"
          value = var.postgres_password
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [google_project_service.run]
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.this.location
  project  = google_cloud_run_service.this.project
  service  = google_cloud_run_service.this.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
