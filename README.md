# goquery namespace cache validation

- Candidate: `788c2d11b29f2d27beeecebb12d1604fcc79ec94`; baseline: `e265bec92fa6513b99a3af99af4830ae3149f282`.
- Based on [PuerkitoBio/goquery#618](https://github.com/PuerkitoBio/goquery/pull/618), original commit `8ece687ac2ec65713939d6fd7a80e80cb836abda`. Shaisolaris's author/date and cherry-pick source are retained in `c2c175d8d083f3fb5d26318b5e03efe9827a08e7`.
- Increment: one 64-line regression function, eight public insertion/wrapping methods in two selection orders, with HTML/SVG/MathML contexts and repeated SVG contexts. No other bug, API, dependency or workflow change.

## Results

The exact upstream `go test ./... -v -cover` command passed locally on Windows and Linux with Go 1.26.10 and 1.27.3. Baseline has 457 PASS records; candidate has 476, including parent tests and subtests (not 476 independent tests). Coverage is 88.4%. Both versions/OS pass vet and build; Linux 1.27 race testing passes.

The same candidate namespace tests on the baseline fail: 18 independent cases (two original tests plus 16 added subtests), producing 19 FAIL events including their parent. The separate executable application imports genuine locally generated Go module ZIPs, asserts that no replacement module was used, parses a complete document and checks DOM/text via public methods. In every Windows/Linux and Go-version combination, baseline fails all six cases and candidate passes six cases/18 checks.

Before PR creation, the unchanged upstream workflow also passed all six jobs on the own public fork: Go 1.26.x/1.27.x × Ubuntu/macOS/Windows. [Exact-head run 40959170934](https://github.com/jipeng6036-del/goquery/actions/runs/40959170934). These fork results are separate from the upstream PR's later CI and maintainer review.

## Reproduce project checks

Checkout the candidate and run, with the desired Go version in PATH:

```sh
go version
go test ./... -v -cover
go vet ./...
go build ./...
gofmt -l manipulation.go manipulation_test.go
go test -race ./...  # supported Linux toolchain
```

For the negative control, create a separate baseline checkout at the SHA above and copy only the candidate `manipulation_test.go` into it. Run `go test -count=1 -v -run NamespaceContext ./...`; failure is expected. Do not apply the production fix to this control.

## Reproduce the installed module application

The two `canonical-proxy-*` directories are file-based Go module proxies containing exact source ZIPs and their `.mod`/`.info`/`list` metadata. `canonical-consumer-*` contains the exact application and input used. The fixture version `v0.0.0-canonical` is local only and was never published as a release. Use an independent module cache for each stage so the shared fixture version cannot reuse the other ZIP.

POSIX, from this packet directory, with an installed Go 1.26.10 or 1.27.3:

```sh
proof="$(pwd)"
for stage in baseline candidate; do
  (cd "canonical-consumer-$stage" && \
   GOTOOLCHAIN=local GOSUMDB=off \
   GOMODCACHE="$proof/repro-cache-$stage" \
   GOPROXY="file://$proof/canonical-proxy-$stage,https://proxy.golang.org" go run .)
done
```

PowerShell, from this packet directory:

```powershell
$proof = (Get-Location).Path
foreach ($stage in @('baseline', 'candidate')) {
  $env:GOTOOLCHAIN = 'local'
  $env:GOSUMDB = 'off'
  $env:GOMODCACHE = Join-Path $proof "repro-cache-$stage"
  $proxy = ([Uri](Join-Path $proof "canonical-proxy-$stage")).AbsoluteUri
  $env:GOPROXY = "$proxy,https://proxy.golang.org"
  Push-Location "canonical-consumer-$stage"
  try { go run . } finally { Pop-Location }
}
```

Baseline exits 1 with six failures; candidate exits 0. `GOSUMDB=off` is limited to these unpublished local fixtures. Run this in a disposable shell, use the exact installed Go binary and do not change production settings.

## Evidence boundaries

- Manifests bind all 64 source file bytes to Git blobs, with original Git modes recorded. Original Git/ZIP has `misc/git/pre-commit` executable; the Windows-extracted tar used by local Linux testing did not retain that unused helper's executable mode. No full tar mode identity is claimed. Fork CI uses normal Git checkout.
- Whole-repository gofmt lists only the identical unchanged baseline/candidate `doc.go` comment formatting. Both changed files are clean; upstream CI has no format gate. This change does not format unrelated files.
- Earlier CRLF exports and a Go toolchain-bootstrap diagnostic are retained locally but are not used as final passing evidence. All reported final checks use canonical Git-blob source and verified actual runtimes.
- `COMMANDS.json` selects the actual canonical command receipts. `$PROOF`, `$CHECKOUT`, `$GO126`, `$GO127`, `$GO126_HOME` and `$GO127_HOME` replace only local paths/toolchain prefixes; these are evidence placeholders, not literal runnable commands. Original evidence is unchanged locally; `PROVENANCE.json` distinguishes original and public SHA-256 values.
- This is code/runtime validation, not a claim of human review, original-author endorsement, paid service validation or upstream approval. AI assistance was used.
