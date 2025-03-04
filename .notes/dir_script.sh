#!/usr/bin/env bash
#
# dir_script.sh
#
# Recursively gather the directory structure and output it as
# a Markdown file, focusing on user-written code and excluding
# dependencies and build artifacts.

################################################################################
# Configuration
################################################################################

PROJECT_ROOT="$(pwd)"           # Use current directory as project root
OUTPUT_FILE="$(pwd)/.notes/directory_structure.md"  # Output file path
MAX_DEPTH=4                     # Maximum depth to traverse
MAX_FILES_PER_DIR=15            # Maximum number of files to show per directory

# Directories to exclude
EXCLUDE_DIRS=(
  # Version control
  ".git" ".github" 

  # Build artifacts and dependencies
  "node_modules" "dist" "release" "build"
  "renderer/node_modules" "renderer/build" "renderer/dist"
  
  # Test directories
  "test" "tests" "coverage" "cypress" "testutil"
  
  # Logs and temporary files
  "logs" "tmp" ".elixir_ls" ".vscode" ".idea"
  
  # Generated documentation
  "doc" "docs" "apidoc"
  
  # Other common non-user code directories
  "out" "target" "bin" "obj"
  
  # Duplicate directories to fix the issue
  ".notes/.notes"
)

# File patterns to exclude
EXCLUDE_FILE_PATTERNS=(
  "*.beam" "*.o" "*.so" "*.dll" "*.exe" "*.lock"
  "*.log" "*.tmp" "*.temp" "*.bak" "*.backup"
  "*.min.js" "*.min.css" "*.map" "*.chunk.js"
  "package-lock.json" "yarn.lock"
  "*.gz" "*.zip" "*.tar" "*.rar" "*.7z"
  "*.png" "*.jpg" "*.jpeg" "*.gif" "*.svg" "*.ico"
  "*.ttf" "*.woff" "*.woff2" "*.eot" ".env.local" "*.csv"
)

# Directories to focus on
FOCUS_DIRS=(
  "internal" 
  "renderer"
  "config"
  "scripts"
  "assets"
)

# Files that are allowed at the root level
ROOT_ALLOWED_FILES=(
  "main.go"
  "README.md"
  "main.js"
  "go.mod"
  "go.sum"
  "preload.js"
  "package.json"
  "postcss.config.js"
  "tailwind.config.js"
  "eslint.config.mjs"
  ".gitignore"
  "CHANGELOG.md"
  ".env"
  "version"
)

################################################################################
# Helper to build the indent prefix
################################################################################

indent_prefix() {
  local indent="$1"
  local prefix=""
  for ((i=0; i<indent; i++)); do
    prefix+="    "
  done
  echo -n "$prefix"
}

################################################################################
# Helper to check if path should be excluded
################################################################################

should_exclude() {
  local path="$1"
  local rel_path="${path#$PROJECT_ROOT/}"  # Remove PROJECT_ROOT prefix
  local indent="$2"
  
  # For root level, only include specific files (whitelist approach)
  if [[ $indent -eq 0 && -f "$path" ]]; then
    local basename=$(basename "$path")
    local allowed=0
    for allow_file in "${ROOT_ALLOWED_FILES[@]}"; do
      if [[ "$basename" == "$allow_file" ]]; then
        allowed=1
        break
      fi
    done
    
    if [[ $allowed -eq 0 ]]; then
      return 0  # true, should exclude
    fi
  fi
  
  # Check if path is in focus directories
  local in_focus=0
  for focus in "${FOCUS_DIRS[@]}"; do
    if [[ "$rel_path" == "$focus" || "$rel_path" == "$focus/"* ]]; then
      in_focus=1
      break
    fi
  done
  
  # If not in focus directories, exclude
  if [[ $in_focus -eq 0 ]]; then
    return 0  # true, should exclude
  fi
  
  # Check directory exclusions
  for exclude in "${EXCLUDE_DIRS[@]}"; do
    if [[ "$rel_path" == "$exclude" || "$rel_path" == "$exclude/"* ]]; then
      return 0  # true, should exclude
    fi
  done
  
  # If it's a file, check file pattern exclusions
  if [[ -f "$path" ]]; then
    for pattern in "${EXCLUDE_FILE_PATTERNS[@]}"; do
      if [[ "$path" == *$pattern ]]; then
        return 0  # true, should exclude
      fi
    done
  fi
  
  return 1  # false, should not exclude
}

################################################################################
# Helper to get file type icon
################################################################################

get_file_icon() {
  local file="$1"
  local extension="${file##*.}"
  
  case "$extension" in
    go)          echo "🔹" ;; # Go
    js|jsx)      echo "🟨" ;; # JavaScript
    ts|tsx)      echo "🔵" ;; # TypeScript
    css|scss)    echo "🎨" ;; # Stylesheets
    json)        echo "📋" ;; # JSON
    md)          echo "📝" ;; # Markdown
    html|htm)    echo "🌐" ;; # HTML
    sql)         echo "💾" ;; # SQL
    sh|bash)     echo "🐚" ;; # Shell scripts
    yml|yaml)    echo "⚙️" ;; # YAML
    mod|sum)     echo "📦" ;; # Go modules
    cjs|mjs)     echo "🟧" ;; # CommonJS/ES modules
    *)           echo "📄" ;; # Other files
  esac
}

################################################################################
# Helper to get directory type icon
################################################################################

get_dir_icon() {
  local dir="$1"
  local basename=$(basename "$dir")
  
  case "$basename" in
    internal)    echo "🔹" ;; # Go backend code
    handlers)    echo "🎮" ;; # API handlers
    services)    echo "⚙️" ;; # Services
    model)       echo "💎" ;; # Models
    persist)     echo "💾" ;; # Persistence
    cmd)         echo "🚀" ;; # Commands
    server)      echo "🖥️" ;; # Server
    http)        echo "🌐" ;; # HTTP
    errors)      echo "⚠️" ;; # Errors
    embed)       echo "📦" ;; # Embedded resources
    renderer)    echo "🎭" ;; # Frontend renderer
    src)         echo "📁" ;; # Source code
    components)  echo "🧩" ;; # Components
    hooks)       echo "🪝" ;; # Hooks
    pages)       echo "📄" ;; # Pages
    assets)      echo "🖼️" ;; # Assets
    api)         echo "🔌" ;; # API
    utils)       echo "🔧" ;; # Utilities
    config)      echo "⚙️" ;; # Configuration
    scripts)     echo "📜" ;; # Scripts
    public)      echo "🌍" ;; # Public assets
    .notes)      echo "📔" ;; # Notes
    *)           echo "📁" ;; # Other directories
  esac
}

################################################################################
# Recursive function to list files and directories with Markdown formatting
################################################################################

get_formatted_directory() {
  local path="$1"
  local indent="${2:-0}"
  local max_depth="${3:-$MAX_DEPTH}"  # Use global MAX_DEPTH as default
  
  # Stop if we've reached max depth
  if [[ $indent -ge $max_depth ]]; then
    local prefix
    prefix="$(indent_prefix "$indent")"
    echo "${prefix}- *(max depth reached)*"
    return
  fi

  # Read directory contents into an array, ignoring errors
  local items=()
  IFS=$'\n' read -r -d '' -a items < <(ls -A "$path" 2>/dev/null && printf '\0')
  
  # Sort items: directories first, then files
  local dirs=()
  local files=()
  
  for item in "${items[@]}"; do
    # Ignore the current and parent dir entries
    [[ "$item" == "." || "$item" == ".." ]] && continue
    
    local fullpath="$path/$item"
    
    # Skip excluded directories and files
    should_exclude "$fullpath" $indent && continue
    
    if [[ -d "$fullpath" ]]; then
      dirs+=("$item")
    else
      # Only add files if not at root level
      if [[ $indent -gt 0 ]]; then
        files+=("$item")
      fi
    fi
  done
  
  # Sort directories and files alphabetically
  IFS=$'\n' sorted_dirs=($(sort <<<"${dirs[*]}"))
  IFS=$'\n' sorted_files=($(sort <<<"${files[*]}"))
  
  # Process directories
  for item in "${sorted_dirs[@]}"; do
    local fullpath="$path/$item"
    local prefix
    prefix="$(indent_prefix "$indent")"
    local icon=$(get_dir_icon "$fullpath")
    
    echo "${prefix}- ${icon} **${item}/**"
    get_formatted_directory "$fullpath" $((indent + 1)) $max_depth
  done
  
  # Process files (limited to MAX_FILES_PER_DIR) - only if not at root level
  if [[ $indent -gt 0 ]]; then
    local file_count=0
    for item in "${sorted_files[@]}"; do
      # Limit the number of files shown per directory
      if [[ $file_count -ge $MAX_FILES_PER_DIR ]]; then
        local prefix
        prefix="$(indent_prefix "$indent")"
        echo "${prefix}- *(and $(( ${#sorted_files[@]} - $MAX_FILES_PER_DIR )) more files)*"
        break
      fi
      
      local fullpath="$path/$item"
      local prefix
      prefix="$(indent_prefix "$indent")"
      local icon=$(get_file_icon "$item")
      
      echo "${prefix}- ${icon} ${item}"
      file_count=$((file_count + 1))
    done
  fi
}

################################################################################
# Generate the output content (Markdown)
################################################################################

# Store in a variable so we can write once at the end
MARKDOWN_CONTENT="# CanIFly Project Directory Structure

This document provides an overview of the project's directory structure, focusing on user-written code and excluding dependencies and build artifacts.

## Legend
- 🔹 Go (*.go)
- 🟨 JavaScript (*.js, *.jsx)
- 🔵 TypeScript (*.ts, *.tsx)
- 🎨 Stylesheets (*.css, *.scss)
- 📋 JSON (*.json)
- 📝 Markdown (*.md)
- 🌐 HTML (*.html, *.htm)
- 💾 SQL (*.sql)
- 🐚 Shell Scripts (*.sh, *.bash)
- ⚙️ Configuration (*.yml, *.yaml, config)
- 📦 Go modules (go.mod, go.sum)
- 🟧 CommonJS/ES modules (*.cjs, *.mjs)
- 📄 Other Files

## Core Components

\`\`\`
$( get_formatted_directory "$PROJECT_ROOT" 0 $MAX_DEPTH )
\`\`\`

## Note
This structure was automatically generated and may not include all files. Directories and files that are typically not user code (build artifacts, dependencies, etc.) have been excluded. The structure is limited to a depth of $MAX_DEPTH levels and shows at most $MAX_FILES_PER_DIR files per directory.
"

################################################################################
# Write to the output file
################################################################################

# Ensure output directory exists
mkdir -p "$(dirname "$OUTPUT_FILE")"
# Write content
echo "$MARKDOWN_CONTENT" > "$OUTPUT_FILE"

echo "Directory structure updated in $OUTPUT_FILE"

