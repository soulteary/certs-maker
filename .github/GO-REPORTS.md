# Go report workflow

The PR job runs both owned Actions with `commit: false` and released CLIs
(`go-test-report-action v1.1.2`, `goreportcard-action v1.1.3`, pinned wrapper SHAs).
It replaces the primary Linux Go test/coverage invocation and reuses its profile
for existing coverage artifacts/gates. Other OS matrices, release verification,
integration, benchmarks and fuzz workflows retain their distinct checks.

Reports run on pushed tags and `workflow_dispatch`. They do not generate bot
commits on branch pushes. A separate job may commit the five stable reports in
one report commit, only if the tested commit still equals the default branch
head. Older release tags and non-default manual refs retain downloadable reports
without overwriting the current branch. PRs never enter that write job. Protected
branches can reject writeback; no PAT, bypass, or force-push is used.

Test directory: `.`. Packages: `./...`.
Race detector: `true`. Additional args: `-count=1`.
Coverage floor: `0%` (zero records a baseline).
Per-package floor: `0%`.
Quality-score floor: `0%`.
Raw profiles/test logs are artifacts, never committed.

Published releases query actual tag-push report runs for the tested commit. A successful or active tag run prevents duplicate tests; missing, failed, cancelled, or skipped runs allow the release event to generate reports. This covers skipped commits and releases whose tag push produced no workflow event. API lookup errors fail the gate rather than silently assume a report exists. Stable report-only artifacts explicitly restore `.github` on download; mixed report/document artifacts retain repository hierarchy.

Before generation the job removes only its known generated reports and raw result directory, so failed producers cannot upload stale checked-in evidence. Artifact names include the run attempt. The write token already prevents recursive workflow triggers, so report commits omit CI skip annotations that could suppress a later release tag.
