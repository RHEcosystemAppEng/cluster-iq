#!/bin/bash
set -e

# Set default backend URL if not provided
BACKEND_URL="${BACKEND_URL:-http://localhost:8081}"
export BACKEND_URL

echo "INFO: Using BACKEND_URL=${BACKEND_URL}"

# Paths to the template and generated NGINX configuration
NGINX_TEMPLATE="/etc/nginx/nginx.conf.template"
NGINX_CONFIG="/etc/nginx/nginx.conf"

# Ensure the template file exists
if [[ ! -f "$NGINX_TEMPLATE" ]]; then
    echo "ERROR: Template file '${NGINX_TEMPLATE}' not found."
    exit 1
fi

# Replace environment variables in the NGINX configuration
echo "INFO: Substituting environment variables in NGINX template..."
envsubst '${BACKEND_URL}' < "$NGINX_TEMPLATE" > "$NGINX_CONFIG"

# Validate if the substitution was successful
if grep -q '${BACKEND_URL}' "$NGINX_CONFIG"; then
    echo "ERROR: Environment variable substitution failed in '${NGINX_CONFIG}'."
    exit 1
fi

# Display the generated configuration
echo "INFO: NGINX configuration generated successfully:"
grep "proxy_pass" "$NGINX_CONFIG"
