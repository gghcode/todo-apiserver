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
          name  = "BASE_PATH"
          value = var.base_path
        }
        env {
          name  = "BCRYPT_COST"
          value = var.bcrypt_cost
        }
        env {
          name  = "GRACEFUL_SHUTDOWN_TIMEOUT_SEC"
          value = var.graceful_shutdown_timeout_sec
        }
        env {
          name  = "JWT_SECRET_KEY"
          value = var.jwt_secret_key
        }
        env {
          name  = "JWT_ACCESS_EXPIRES_IN_SEC"
          value = var.jwt_access_expires_in_sec
        }
        env {
          name  = "JWT_REFRESH_EXPIRES_IN_SEC"
          value = var.jwt_refresh_expires_in_sec
        }
        env {
          name  = "CORS_ALLOW_ORIGINS"
          value = var.cors_allow_origins
        }
        env {
          name  = "CORS_ALLOW_METHODS"
          value = var.cors_allow_methods
        }
        env {
          name  = "POSTGRES_HOST"
          value = var.postgres_host
        }
        env {
          name  = "POSTGRES_USER"
          value = var.postgres_user
        }
        env {
          name  = "POSTGRES_NAME"
          value = var.postgres_name
        }
        env {
          name  = "REST_POSTGRES_PASSWORD"
          value = var.postgres_password
        }
        env {
          name  = "REDIS_ADDR"
          value = var.redis_addr
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
