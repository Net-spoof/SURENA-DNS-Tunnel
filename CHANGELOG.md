# Changelog

## DNSFabric Subscriber Portal · 0.8.2.2 (preview)

- Browser-aware `/s/TOKEN`: browsers see the dashboard while subscription clients keep the previous raw encoded data and headers.
- Made the user page Persian RTL with a copy-link control.
- Kept legacy `/sub/TOKEN` strictly raw; raw format query parameters always override browser negotiation.

## DNSFabric Subscriber Portal · 0.8.2.1

- Built an independently version-locked WhiteDNS 0.8.2 server patch at the pinned upstream commit, under the upstream Apache-2.0 terms.
- Added `/s/TOKEN` as a short subscription alias and `/u/TOKEN` as a private usage, expiry and QR dashboard.
- Preserved `/sub/TOKEN` response bytes, subscription tokens, management login path, tunnel accounts and quota history.
- Added a portal button and short subscription link in the existing management UI.
- Published automated amd64/arm64 binaries and a guarded in-place updater with SHA-256 verification, private backup and rollback.
- Added Persian RTL and English rollout guides and upstream compatibility tests.

## Unreleased — DNSFabric branding

- Published detailed, example-driven Persian RTL and English installation manuals (`docs/INSTALLATION.fa.md` and `.en.md`).
- Published multi-VPS DNS delegation and naming tutorials for `dn01`/`dn02`/`dn03` in both languages.
- Added copy-friendly Cloudflare record tables, exact installer inputs and tunnel FQDN length examples to the main README.

- Rebranded the installer and CLI as **DNSFabric**.
- Redesigned the bilingual README with RTL Persian content and LTR code blocks.
- Added dedicated Persian troubleshooting, security and contributing guidance.
- Kept official WhiteDNS v0.8.2 installation pinned.
- Preserved old repository-path compatibility in the downloaded bootstrap script.

## 1.0.0 — Initial development baseline

- Guided fresh-VPS setup for WhiteDNS Server Management.
- Configuration collection and DNS delegation instructions.
- Official checksum-verified upstream bundle retrieval.
- HTTPS and private management URL initialization.
- Post-install status and TLS verification.
- Offline installer unit tests and GitHub Actions integration.

**Verification scope:** CI tests pass for offline checks; fresh-VPS end-to-end field validation remains separately required. No new upstream protocol implementation is claimed.
