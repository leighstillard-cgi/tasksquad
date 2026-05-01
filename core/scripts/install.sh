#!/usr/bin/env bash
#
# install.sh - TaskSquad Environment Bootstrap
#
# Installs and configures the Claude Code plugins and tools that TaskSquad
# depends on. Idempotent - safe to run multiple times.
#
# Usage: ./core/scripts/install.sh [--skip-rtk] [--skip-graphify] [--skip-lasso]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Parse arguments
SKIP_RTK=false
SKIP_GRAPHIFY=false
SKIP_LASSO=false

for arg in "$@"; do
  case "$arg" in
    --skip-rtk) SKIP_RTK=true ;;
    --skip-graphify) SKIP_GRAPHIFY=true ;;
    --skip-lasso) SKIP_LASSO=true ;;
    --help|-h)
      echo "Usage: $0 [--skip-rtk] [--skip-graphify] [--skip-lasso]"
      echo ""
      echo "Options:"
      echo "  --skip-rtk       Skip RTK installation (requires Rust)"
      echo "  --skip-graphify  Skip graphify dependencies check"
      echo "  --skip-lasso     Skip lasso-security hooks (optional)"
      exit 0
      ;;
  esac
done

# Tracking arrays for summary
declare -a INSTALLED=()
declare -a ALREADY_PRESENT=()
declare -a SKIPPED=()
declare -a FAILED=()

# Color output (degrades gracefully)
if [[ -t 1 ]] && command -v tput &>/dev/null; then
  RED=$(tput setaf 1 2>/dev/null || echo "")
  GREEN=$(tput setaf 2 2>/dev/null || echo "")
  YELLOW=$(tput setaf 3 2>/dev/null || echo "")
  BLUE=$(tput setaf 4 2>/dev/null || echo "")
  RESET=$(tput sgr0 2>/dev/null || echo "")
else
  RED="" GREEN="" YELLOW="" BLUE="" RESET=""
fi

info()    { echo "${BLUE}[INFO]${RESET} $*"; }
success() { echo "${GREEN}[OK]${RESET} $*"; }
warn()    { echo "${YELLOW}[WARN]${RESET} $*"; }
error()   { echo "${RED}[ERROR]${RESET} $*"; }

# ============================================================================
# Prerequisites Check
# ============================================================================

info "Checking prerequisites..."

PREREQ_FAIL=false

# Node.js (required for Claude Code and plugins)
if command -v node &>/dev/null; then
  NODE_VERSION=$(node --version)
  success "Node.js: $NODE_VERSION"
else
  error "Node.js is required but not installed."
  error "Install via: https://nodejs.org/ or 'nvm install --lts'"
  PREREQ_FAIL=true
fi

# Python 3 (required for graphify and wiki linting)
if command -v python3 &>/dev/null; then
  PYTHON_VERSION=$(python3 --version)
  success "Python: $PYTHON_VERSION"
else
  error "Python 3 is required but not installed."
  PREREQ_FAIL=true
fi

# Claude CLI (required)
if command -v claude &>/dev/null; then
  success "Claude CLI: found"
else
  error "Claude CLI is required but not installed."
  error "Install via: npm install -g @anthropic-ai/claude-cli"
  PREREQ_FAIL=true
fi

# Rust (optional - needed for RTK)
if command -v cargo &>/dev/null; then
  CARGO_VERSION=$(cargo --version)
  success "Rust/Cargo: $CARGO_VERSION"
  RUST_AVAILABLE=true
else
  warn "Rust/Cargo not installed - RTK will be skipped"
  warn "Install via: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
  RUST_AVAILABLE=false
fi

# Git (required)
if command -v git &>/dev/null; then
  GIT_VERSION=$(git --version)
  success "Git: $GIT_VERSION"
else
  error "Git is required but not installed."
  PREREQ_FAIL=true
fi

if [[ "$PREREQ_FAIL" == "true" ]]; then
  error ""
  error "Prerequisites not met. Install the missing tools and re-run."
  exit 1
fi

echo ""

# ============================================================================
# Plugin Installation
# ============================================================================

info "Installing Claude Code plugins..."

# claude-mem plugin
if claude plugin list 2>/dev/null | grep -q "claude-mem@thedotmack"; then
  success "claude-mem: already installed"
  ALREADY_PRESENT+=("claude-mem plugin")
else
  info "Installing claude-mem plugin..."
  if claude plugin install claude-mem@thedotmack --scope user 2>/dev/null; then
    success "claude-mem: installed"
    INSTALLED+=("claude-mem plugin")
  else
    warn "claude-mem: failed to install (may need marketplace added first)"
    FAILED+=("claude-mem plugin")
  fi
fi

echo ""

# ============================================================================
# Graphify Dependencies
# ============================================================================

if [[ "$SKIP_GRAPHIFY" == "true" ]]; then
  info "Skipping graphify (--skip-graphify)"
  SKIPPED+=("graphify dependencies")
else
  info "Checking graphify dependencies..."

  # Check if graphify Python package is installed
  if python3 -c "import graphify" 2>/dev/null; then
    success "graphify Python package: already installed"
    ALREADY_PRESENT+=("graphify Python package")
  else
    info "Installing graphify Python package..."
    if python3 -m pip install graphifyy -q 2>/dev/null || \
       python3 -m pip install graphifyy -q --user 2>/dev/null; then
      success "graphify Python package: installed"
      INSTALLED+=("graphify Python package")
    else
      warn "graphify Python package: failed to install"
      warn "Try manually: python3 -m pip install graphifyy"
      FAILED+=("graphify Python package")
    fi
  fi

  # Check graphify skill exists
  if [[ -f "$HOME/.claude/skills/graphify/SKILL.md" ]]; then
    success "graphify skill: present at ~/.claude/skills/graphify/"
    ALREADY_PRESENT+=("graphify skill")
  else
    warn "graphify skill: not found at ~/.claude/skills/graphify/"
    warn "Copy the skill folder manually or install from marketplace"
    FAILED+=("graphify skill")
  fi
fi

echo ""

# ============================================================================
# RTK Installation
# ============================================================================

if [[ "$SKIP_RTK" == "true" ]]; then
  info "Skipping RTK (--skip-rtk)"
  SKIPPED+=("RTK")
elif [[ "$RUST_AVAILABLE" == "false" ]]; then
  info "Skipping RTK (Rust not available)"
  SKIPPED+=("RTK (no Rust)")
else
  info "Checking RTK installation..."

  if command -v rtk &>/dev/null; then
    RTK_VERSION=$(rtk --version 2>/dev/null || echo "unknown")
    success "RTK: already installed ($RTK_VERSION)"
    ALREADY_PRESENT+=("RTK")
  else
    info "Installing RTK via cargo..."
    if cargo install rtk 2>/dev/null; then
      success "RTK: installed"
      INSTALLED+=("RTK")
    else
      warn "RTK: failed to install"
      warn "Try manually: cargo install rtk"
      FAILED+=("RTK")
    fi
  fi
fi

echo ""

# ============================================================================
# .NET LSP Plugin (Conditional)
# ============================================================================

info "Checking .NET SDK..."

if command -v dotnet &>/dev/null; then
  DOTNET_VERSION=$(dotnet --version 2>/dev/null || echo "unknown")
  success ".NET SDK: $DOTNET_VERSION"

  # Check if csharp-ls LSP plugin is enabled
  if claude plugin list 2>/dev/null | grep -q "csharp-lsp@claude-plugins-official"; then
    if claude plugin list 2>/dev/null | grep -A1 "csharp-lsp@claude-plugins-official" | grep -q "enabled"; then
      success "C# LSP plugin: already enabled"
      ALREADY_PRESENT+=("C# LSP plugin")
    else
      info "Enabling C# LSP plugin..."
      # Plugin exists but is disabled - would need settings.json update
      warn "C# LSP plugin: found but disabled. Enable in settings.json"
      SKIPPED+=("C# LSP plugin (needs manual enable)")
    fi
  else
    info "C# LSP plugin not installed. Install via: claude plugin install csharp-lsp@claude-plugins-official"
    SKIPPED+=("C# LSP plugin (not installed)")
  fi
else
  info ".NET SDK not detected - skipping C# LSP configuration"
  SKIPPED+=(".NET LSP (no SDK)")
fi

echo ""

# ============================================================================
# Settings.json Configuration
# ============================================================================

info "Checking Claude Code settings..."

SETTINGS_FILE="$HOME/.claude/settings.json"

if [[ -f "$SETTINGS_FILE" ]]; then
  success "settings.json: exists at $SETTINGS_FILE"

  # Check for thedotmack marketplace
  if grep -q '"thedotmack"' "$SETTINGS_FILE" 2>/dev/null; then
    success "thedotmack marketplace: configured"
    ALREADY_PRESENT+=("thedotmack marketplace config")
  else
    warn "thedotmack marketplace: not configured in settings.json"
    warn "Add to extraKnownMarketplaces for claude-mem plugin"
    SKIPPED+=("thedotmack marketplace (manual add needed)")
  fi

  # Check for claude-mem enabled
  if grep -q '"claude-mem@thedotmack": true' "$SETTINGS_FILE" 2>/dev/null; then
    success "claude-mem plugin: enabled in settings"
    ALREADY_PRESENT+=("claude-mem enabled in settings")
  else
    warn "claude-mem plugin: not enabled in settings.json"
    SKIPPED+=("claude-mem enable (manual add needed)")
  fi
else
  warn "settings.json: not found at $SETTINGS_FILE"
  warn "Run 'claude' once to generate default settings, then re-run this script"
  FAILED+=("settings.json check")
fi

echo ""

# ============================================================================
# Lasso Security Hooks (Optional)
# ============================================================================

if [[ "$SKIP_LASSO" == "true" ]]; then
  info "Skipping lasso-security hooks (--skip-lasso)"
  SKIPPED+=("lasso-security hooks")
else
  info "Checking lasso-security hooks (optional)..."

  LASSO_DIR="$HOME/.claude/lasso-hooks"
  if [[ -d "$LASSO_DIR" ]]; then
    success "lasso-security hooks: found at $LASSO_DIR"
    ALREADY_PRESENT+=("lasso-security hooks")
  else
    info "lasso-security hooks: not installed (optional)"
    info "Install from: https://github.com/lasso-security/claude-hooks"
    SKIPPED+=("lasso-security hooks (optional)")
  fi
fi

echo ""

# ============================================================================
# Summary
# ============================================================================

echo "========================================"
echo "         INSTALLATION SUMMARY"
echo "========================================"
echo ""

if [[ ${#INSTALLED[@]} -gt 0 ]]; then
  echo "${GREEN}INSTALLED:${RESET}"
  for item in "${INSTALLED[@]}"; do
    echo "  + $item"
  done
  echo ""
fi

if [[ ${#ALREADY_PRESENT[@]} -gt 0 ]]; then
  echo "${BLUE}ALREADY PRESENT:${RESET}"
  for item in "${ALREADY_PRESENT[@]}"; do
    echo "  = $item"
  done
  echo ""
fi

if [[ ${#SKIPPED[@]} -gt 0 ]]; then
  echo "${YELLOW}SKIPPED:${RESET}"
  for item in "${SKIPPED[@]}"; do
    echo "  - $item"
  done
  echo ""
fi

if [[ ${#FAILED[@]} -gt 0 ]]; then
  echo "${RED}FAILED:${RESET}"
  for item in "${FAILED[@]}"; do
    echo "  ! $item"
  done
  echo ""
fi

# ============================================================================
# Write SETUP_COMPLETE.md
# ============================================================================

SETUP_COMPLETE="$REPO_ROOT/SETUP_COMPLETE.md"
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

cat > "$SETUP_COMPLETE" << EOF
# TaskSquad Setup Complete

**Generated:** $TIMESTAMP
**Script:** core/scripts/install.sh

## Installation Results

### Installed
$(if [[ ${#INSTALLED[@]} -gt 0 ]]; then for item in "${INSTALLED[@]}"; do echo "- $item"; done; else echo "- (none)"; fi)

### Already Present
$(if [[ ${#ALREADY_PRESENT[@]} -gt 0 ]]; then for item in "${ALREADY_PRESENT[@]}"; do echo "- $item"; done; else echo "- (none)"; fi)

### Skipped
$(if [[ ${#SKIPPED[@]} -gt 0 ]]; then for item in "${SKIPPED[@]}"; do echo "- $item"; done; else echo "- (none)"; fi)

### Failed
$(if [[ ${#FAILED[@]} -gt 0 ]]; then for item in "${FAILED[@]}"; do echo "- $item"; done; else echo "- (none)"; fi)

## Post-Setup Checklist

Complete these steps after installation:

- [ ] **Rebuild graphify knowledge graph** - Run \`/graphify\` on your content directories
- [ ] **Run wiki lint** - Execute \`./core/scripts/lint-wiki.sh\` to check wiki structure
- [ ] **Populate canonical-facts.md** - Add project-specific facts to \`data/project/data/canonical-facts.md\`
- [ ] **Configure lasso-security hooks** (optional) - Clone from https://github.com/lasso-security/claude-hooks

## Manual Configuration (if needed)

If any plugins failed to install automatically, add this to \`~/.claude/settings.json\`:

\`\`\`json
{
  "extraKnownMarketplaces": {
    "thedotmack": {
      "source": {
        "source": "github",
        "repo": "thedotmack/claude-mem"
      }
    }
  },
  "enabledPlugins": {
    "claude-mem@thedotmack": true
  }
}
\`\`\`

## Next Steps

1. Run \`./core/scripts/post-setup.sh\` after adding your content
2. Start a Claude Code session and verify plugins are working
3. Test with \`/graphify\` to build your knowledge graph
EOF

success "Setup complete file written to: $SETUP_COMPLETE"

echo ""
echo "========================================"
echo "         POST-SETUP CHECKLIST"
echo "========================================"
echo ""
echo "1. Rebuild graphify knowledge graph:"
echo "   /graphify <your-content-directory>"
echo ""
echo "2. Run wiki lint:"
echo "   ./core/scripts/lint-wiki.sh"
echo ""
echo "3. Populate canonical-facts.md:"
echo "   Edit data/project/data/canonical-facts.md with project facts"
echo ""
echo "4. (Optional) Install lasso-security hooks:"
echo "   https://github.com/lasso-security/claude-hooks"
echo ""
echo "Run ./core/scripts/post-setup.sh after adding content to rebuild indexes."
echo ""

# Exit with error if any critical components failed
if [[ ${#FAILED[@]} -gt 0 ]]; then
  warn "Some components failed to install. Check the summary above."
  exit 1
fi

success "Installation complete!"
exit 0
