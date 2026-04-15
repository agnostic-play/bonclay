#!/usr/bin/env bash
set -o pipefail

export GOPROXY="${GOPROXY:-direct}"

# === CONFIGURATION ===
NEXUS_URL="http://california.allobank.local:8081"
REPO_NAME="golang-hosted"
NEXUS_USER="bondowoso"
NEXUS_PASS="P@ssw0rdC4nd!"                          # Nexus password

GOMOD_FILE="../go.mod"
GOMODCACHE="$(go env GOMODCACHE)"

# Debug mode (set to 1 to enable)
DEBUG=${DEBUG:-0}

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "[INFO] Starting upload of modules from $GOMOD_FILE to Nexus..."
echo "[INFO] Go module cache: $GOMODCACHE"

# === Go Proxy path escaping (per 'go help goproxy') ===
# Uppercase -> !lowercase ; '!' -> '!!'
proxy_escape() {
  local in="$1" out="" ch
  for ((i=0; i<${#in}; i++)); do
    ch="${in:i:1}"
    if [[ "$ch" == "!" ]]; then
      out+="!!"
    elif [[ "$ch" =~ [A-Z] ]]; then
      out+="!${ch,,}"
    else
      out+="$ch"
    fi
  done
  printf '%s' "$out"
}

# === FUNCTION TO CHECK IF MODULE EXISTS IN NEXUS ===
check_module_exists() {
  local module="$1"      # UNESCAPED
  local version="$2"
  local escaped_module
  escaped_module="$(proxy_escape "$module")"
  local upload_path="${escaped_module}/@v"
  local check_url="$NEXUS_URL/repository/$REPO_NAME/$upload_path/${version}.mod"

  curl -s -f -u "$NEXUS_USER:$NEXUS_PASS" "$check_url" -o /dev/null
}

# === EXTRACT MODULES FROM go.mod ===
extract_modules() {
  awk '
    /^require[[:space:]]+\(/ { in_block = 1; next }
    /^require[[:space:]]+/ && !/\(/ {
        gsub(/\/\/.*$/, "", $0)
        if ($2 && $3) print $2, $3
    }
    in_block && /^[[:space:]]+[^\/)]/ {
        gsub(/\/\/.*$/, "", $0)
        if ($1 !~ /^\/\// && $1 && $2) print $1, $2
    }
    in_block && /^\)/ { in_block = 0 }
  ' "$GOMOD_FILE" | grep -v '^[[:space:]]*$' | sort -u
}

# === MAIN PROCESSING ===
total_modules=0
uploaded_modules=0
skipped_modules=0
failed_modules=0

# Materialize the module list to avoid subshell issues
mapfile -t MODULE_LINES < <(extract_modules)
total_modules="${#MODULE_LINES[@]}"
echo "[INFO] Found $total_modules modules to process"
echo ""

for line in "${MODULE_LINES[@]}"; do
  # Robust split on whitespace
  MODULE="${line%%[[:space:]]*}"
  VERSION="${line#*[[:space:]]}"
  if [[ -z "$MODULE" || -z "$VERSION" || "$MODULE" == "$VERSION" ]]; then
    continue
  fi

  echo -e "\n[INFO] Processing ${MODULE}@${VERSION}"

  # === CHECK IF MODULE ALREADY EXISTS IN NEXUS ===
  if check_module_exists "$MODULE" "$VERSION"; then
    echo -e "${YELLOW}[SKIP] Module ${MODULE}@${VERSION} already exists in Nexus${NC}"
    ((skipped_modules++)) || true
    continue
  fi

  # === DOWNLOAD MODULE ===
  echo "[INFO] Downloading module..."
  if ! GOPROXY=direct go mod download "${MODULE}@${VERSION}" 2>/dev/null; then
    echo -e "${RED}[ERROR] Failed to download ${MODULE}@${VERSION}. Skipping...${NC}"
    ((failed_modules++)) || true
    continue
  fi

  # Cache directory uses PROXY-ESCAPED module path
  ESCAPED_MODULE="$(proxy_escape "$MODULE")"
  MODULE_DIR="${GOMODCACHE}/cache/download/${ESCAPED_MODULE}/@v"

  # Debug: List files in the module directory
  if [[ "$DEBUG" = "1" ]]; then
    echo "[DEBUG] Checking directory: $MODULE_DIR"
    if [[ -d "$MODULE_DIR" ]]; then
      echo "[DEBUG] Files in directory (zip/mod/info):"
      ls -la "$MODULE_DIR" | grep -E "\.(zip|mod|info)$" || echo "[DEBUG] No matching files found"
    else
      echo "[DEBUG] Directory does not exist!"
    fi
  fi

  # Files are generally saved under the EXACT version (including +incompatible)
  ZIP_FILE="${MODULE_DIR}/${VERSION}.zip"
  MOD_FILE="${MODULE_DIR}/${VERSION}.mod"
  INFO_FILE="${MODULE_DIR}/${VERSION}.info"

  # Fallback: some edge cases store without '+incompatible'—try a cleaned variant if missing
  if [[ ! -f "$ZIP_FILE" || ! -f "$MOD_FILE" || ! -f "$INFO_FILE" ]]; then
    CLEAN_VERSION="${VERSION%%+incompatible}"
    if [[ "$CLEAN_VERSION" != "$VERSION" ]]; then
      if [[ -f "${MODULE_DIR}/${CLEAN_VERSION}.zip" ]]; then ZIP_FILE="${MODULE_DIR}/${CLEAN_VERSION}.zip"; fi
      if [[ -f "${MODULE_DIR}/${CLEAN_VERSION}.mod" ]]; then MOD_FILE="${MODULE_DIR}/${CLEAN_VERSION}.mod"; fi
      if [[ -f "${MODULE_DIR}/${CLEAN_VERSION}.info" ]]; then INFO_FILE="${MODULE_DIR}/${CLEAN_VERSION}.info"; fi
      # Keep VERSION for upload naming if the canonical files existed under cleaned variant
      if [[ -f "$ZIP_FILE" && -f "$MOD_FILE" && -f "$INFO_FILE" ]]; then
        echo "[INFO] Using files stored under cleaned variant: ${CLEAN_VERSION}"
        VERSION="${CLEAN_VERSION}"
      fi
    fi
  fi

  # Check if all required files exist
  if [[ ! -f "$ZIP_FILE" || ! -f "$MOD_FILE" || ! -f "$INFO_FILE" ]]; then
    echo -e "${RED}[ERROR] Missing files for ${MODULE}@${VERSION}. Skipping...${NC}"
    echo "  Expected files:"
    echo "    ZIP:  $ZIP_FILE (exists: $([ -f "$ZIP_FILE" ] && echo yes || echo no))"
    echo "    MOD:  $MOD_FILE (exists: $([ -f "$MOD_FILE" ] && echo yes || echo no))"
    echo "    INFO: $INFO_FILE (exists: $([ -f "$INFO_FILE" ] && echo yes || echo no))"
    ((failed_modules++)) || true
    continue
  fi

  # === FORMAT UPLOAD PATH ===
  UPLOAD_PATH="${ESCAPED_MODULE}/@v"
  echo "[INFO] Uploading to Nexus: $UPLOAD_PATH (version: $VERSION)"

  upload_success=true

  # Upload .zip
  if ! curl -s -f -u "$NEXUS_USER:$NEXUS_PASS" \
      --upload-file "$ZIP_FILE" \
      -X PUT \
      "$NEXUS_URL/repository/$REPO_NAME/$UPLOAD_PATH/${VERSION}.zip"; then
    echo -e "${RED}[ERROR] Failed to upload ${VERSION}.zip${NC}"
    upload_success=false
  fi

  # Upload .mod
  if ! curl -s -f -u "$NEXUS_USER:$NEXUS_PASS" \
      --upload-file "$MOD_FILE" \
      -X PUT \
      "$NEXUS_URL/repository/$REPO_NAME/$UPLOAD_PATH/${VERSION}.mod"; then
    echo -e "${RED}[ERROR] Failed to upload ${VERSION}.mod${NC}"
    upload_success=false
  fi

  # Upload .info
  if ! curl -s -f -u "$NEXUS_USER:$NEXUS_PASS" \
      --upload-file "$INFO_FILE" \
      -X PUT \
      "$NEXUS_URL/repository/$REPO_NAME/$UPLOAD_PATH/${VERSION}.info"; then
    echo -e "${RED}[ERROR] Failed to upload ${VERSION}.info${NC}"
    upload_success=false
  fi

  if [[ "$upload_success" == true ]]; then
    echo -e "${GREEN}[SUCCESS] Uploaded ${MODULE}@${VERSION}${NC}"
    ((uploaded_modules++)) || true
  else
    echo -e "${RED}[FAILED] Failed to upload some files for ${MODULE}@${VERSION}${NC}"
    ((failed_modules++)) || true
  fi
done

# === SUMMARY ===
echo -e "\n========================================="
echo "UPLOAD SUMMARY:"
echo "========================================="
echo -e "Total modules:    $total_modules"
echo -e "${GREEN}Uploaded:         $uploaded_modules${NC}"
echo -e "${YELLOW}Skipped (exists): $skipped_modules${NC}"
echo -e "${RED}Failed:           $failed_modules${NC}"
echo "========================================="

if [[ $failed_modules -eq 0 ]]; then
  echo -e "\n${GREEN}✅ All modules processed successfully!${NC}"
else
  echo -e "\n${YELLOW}⚠️  Some modules failed to process. Check the logs above.${NC}"
  exit 1
fi
