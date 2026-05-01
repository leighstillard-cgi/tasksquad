# Troubleshooting

## Permission Denied Running Scripts

Use `bash` explicitly:

```bash
bash core/scripts/install.sh
bash core/scripts/post-setup.sh
```

The scripts are tracked as executable for Unix-like environments, but Windows checkouts can still behave differently depending on Git and filesystem settings.

## Graphify Output Is Missing

`graphify-out/` is generated and ignored by Git. A fresh clone will not include it.

First run:

```bash
bash core/scripts/install.sh
bash core/scripts/post-setup.sh
```

If `graphify-out/GRAPH_REPORT.md` is still missing, open Claude Code in the repository root and run:

```text
/graphify data/wiki core/docs/standards core/templates docs --update
```

If the graphify Python package is missing:

```bash
python3 -m pip install graphifyy
```

If a shell `graphify` command is not available, that is acceptable as long as the Claude Code `/graphify` skill is available.

## Claude Plugin Install Fails

Check that Claude Code is installed and available:

```bash
claude --version
claude plugin list
```

If the plugin marketplace is missing, follow the manual settings guidance printed by `core/scripts/install.sh`.

## Wiki Lint Fails

Run:

```bash
bash core/scripts/lint-wiki.sh
```

Then fix the reported page frontmatter, file naming, or index entry. The linter writes reports under `data/lint-reports/`, which is ignored.

## Dashboard Shows No Data

Confirm these files and directories exist:

- `backlog.md`
- `data/dispatch-log.md`
- `data/completions/`
- `data/dispatches/`
- `data/escalations/`

Run from `core/dashboard` with:

```bash
WORKLOG_PATH=../.. go run .
```

Then open `http://localhost:8080`.

## Bootstrap Health Check Fails

Run:

```bash
bash core/scripts/bootstrap-repo.sh /path/to/target-repo
bash core/scripts/check-repo-health.sh /path/to/target-repo
```

If `CLAUDE.md` already exists in the target repository, back it up or merge its content into `.claude/repo-context.md` before replacing it.

## Hooks Are Not Running

Check:

```bash
ls .claude/hooks
cat .claude/settings.json
```

Claude Code must be started from the repository root for `$PROJECT_DIR` to resolve correctly.

## Windows and WSL Path Issues

Prefer running scripts inside WSL or a devcontainer. If the repository is on the Windows filesystem and scripts behave oddly, clone inside the Linux filesystem for better permissions and line ending behavior.
