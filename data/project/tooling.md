# Project Tooling

> Configuration and usage guides for development tooling.

## Claude-Mem: Cross-Session Memory

claude-mem provides persistent memory across Claude Code sessions. It records observations during work, generates session summaries, and injects recent context when starting new sessions.

### Prerequisites

1. **Worker Service**: claude-mem requires a background worker service running on port 37777
2. **Plugin Installation**: Install via Claude Code plugin marketplace (`thedotmack/claude-mem`)

### Enabling the Plugin

Add to user settings (`~/.claude/settings.json`):

```json
{
  "enabledPlugins": {
    "claude-mem@thedotmack": true
  },
  "extraKnownMarketplaces": {
    "thedotmack": {
      "source": {
        "source": "github",
        "repo": "thedotmack/claude-mem"
      }
    }
  }
}
```

The plugin automatically registers its hooks when enabled. No manual hook configuration is required.

### Hook Chain (Automatic)

claude-mem uses four hooks that fire automatically:

| Hook | Trigger | Purpose |
|------|---------|---------|
| `SessionStart` | Session begins, `/clear`, `/compact` | Starts worker service, injects recent context |
| `UserPromptSubmit` | Each user message | Initializes session tracking |
| `PostToolUse` | After any tool call | Records observations (file reads, edits, bash output) |
| `Stop` | Session ends | Generates and stores session summary |

### Three-Layer Query Workflow

When searching memory, always follow this progression to minimize token usage:

```
1. search(query) → Get index with observation IDs (~50-100 tokens/result)
2. timeline(anchor=ID) → Get chronological context around interesting results
3. get_observations([IDs]) → Fetch full details ONLY for filtered IDs
```

**Never fetch full details without filtering first.** This saves 10x tokens.

#### Layer 1: Search

Find observations matching a keyword or topic:

```
Tool: mcp__plugin_claude-mem_mcp-search__search
Parameters:
  - query: "authentication flow" (required)
  - limit: 10 (default: 20)
  - project: "tasksquad" (optional, filters by project)
  - type: "observation" | "session" | "prompt"
  - dateStart: "2026-04-01"
  - dateEnd: "2026-04-15"
```

Returns a table with observation IDs, timestamps, titles, and estimated token counts.

#### Layer 2: Timeline

Get chronological context around a specific observation:

```
Tool: mcp__plugin_claude-mem_mcp-search__timeline
Parameters:
  - anchor: 12345 (observation ID from search)
  - depth_before: 5 (observations before anchor)
  - depth_after: 5 (observations after anchor)
  - project: "tasksquad" (optional)
```

Use this to understand what happened before/after an interesting observation.

#### Layer 3: Get Observations

Fetch full content for specific observations:

```
Tool: mcp__plugin_claude-mem_mcp-search__get_observations
Parameters:
  - ids: [12345, 12346, 12347] (required, array of IDs)
  - limit: 10
```

Only call this after filtering with search/timeline. Full observations can be large.

### Smart Code Search Tools

claude-mem also provides AST-aware code search:

| Tool | Purpose |
|------|---------|
| `smart_search` | Find symbols, functions, classes across codebase |
| `smart_outline` | Get structural outline of a file (signatures only, bodies folded) |
| `smart_unfold` | Expand a specific symbol to see its full source |

Use `smart_outline` + `smart_unfold` instead of reading entire files to save tokens.

### Verifying Setup

Check worker service health:
```bash
curl http://127.0.0.1:37777/api/health
```

Test observation persistence:
```
# Search for recent activity
mcp__plugin_claude-mem_mcp-search__search query="recent"
```

### Skill Shortcuts

claude-mem provides these skills (invoke via slash command):

- `/claude-mem:mem-search` — Search memory with guided workflow
- `/claude-mem:timeline-report` — Generate project activity report
- `/claude-mem:smart-explore` — Token-optimized code exploration
- `/claude-mem:make-plan` — Create implementation plans with context
- `/claude-mem:do` — Execute phased plans with subagents

### Troubleshooting

**Worker not responding**: The SessionStart hook should auto-start it. If not:
```bash
# Check if process is running
curl http://127.0.0.1:37777/api/health

# Manually start (usually not needed)
node ~/.claude/plugins/cache/thedotmack/claude-mem/*/scripts/worker-service.cjs start
```

**No observations recorded**: Verify the plugin is enabled in settings and the PostToolUse hook is registered. Check Claude Code logs for hook errors.

**Context not injected**: The SessionStart hook injects context into the CLAUDE.md file in the plugin directory. This is then picked up by Claude Code's context loading.
