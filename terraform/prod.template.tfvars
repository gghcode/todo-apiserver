image_tag                     = "${IMAGE_TAG}"
base_path                     = "api"
bcrypt_cost                   = 12
graceful_shutdown_timeout_sec = 30
jwt_access_expires_in_sec     = 3600
jwt_refresh_expires_in_sec    = 86400
cors_allow_origins            = "*"
cors_allow_methods            = "GET,POST,PUT,DELETE,OPTIONS"