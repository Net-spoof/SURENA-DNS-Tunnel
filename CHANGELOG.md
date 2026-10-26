# Changelog

## DNSFabric Subscriber Portal · 0.8.2.7 (preview)

- Added an official direct Android APK download button (WhiteDNS 1.6.14, universal) plus a link to the latest official release page.
- Added a separate iPhone / iOS tile with an Apple-shaped icon, transparently indicating that no verified compatible iOS client release is available.
- No private subscription URLs or account tokens are sent with the APK download; no binaries are repackaged or redistributed.
- Kept all subscription routes, deep-link config imports, quota displays, QR endpoints and service data untouched.
- Added tests for official download URLs, safety of the iOS indication, and download availability on inactive accounts.

## DNSFabric Subscriber Portal · 0.8.2.6 (preview)

- Corrected the intended control: each enabled protocol card now labels its existing direct `soar://`, `cottendns://` or `masterdns://` Open deep link **📱 افزودن به برنامه (Add to App)**.
- Made direct Add to App the visually primary protocol action while leaving **Copy config** separate.
- Renamed the subscription toolbar's copy-and-instructions button to **کپی و راهنمای افزودن اشتراک** to distinguish HTTPS subscription Import from direct profile import.
- Preserved bearer tokens, legacy `/sub/`, smart `/s/`, QR images, quotas and the protected panel login.

## DNSFabric Subscriber Portal · 0.8.2.5 (preview)

- Replaced the misleading browser-facing subscription Import link with **📱 افزودن به برنامه** (Add to App).
- Added a clipboard-copy confirmation and step-by-step in-app subscription Import/Paste dialog, with manual copy fallback.
- Preserved independent **Copy subscription** action, per-engine QR, direct profile Open buttons and unchanged raw `/sub/` for the Android app.
- Documented the current Android client's limitation: HTTPS subscription auto-import is not supported via ACTION_VIEW deep links.

## DNSFabric Subscriber Portal · 0.8.2.4 (preview)

- Responsive redesign of the customer subscription dashboard, with clean two-column desktop/one-column mobile protocol cards.
- Per-engine, per-zone QR endpoint that encodes the exact direct configuration URL, with server-side token, active-status and allowed-zone validation.
- New **Show QR**, QR enlarge modal, and **Open** direct protocol deep-link actions; retained separate Copy Config and Copy Subscription buttons.
- Added regression tests for QR payload matching, invalid zone/engine rejection and disabled-account QR privacy.

## DNSFabric Subscriber Portal · 0.8.2.3 (preview)

- Removed product branding from the public subscriber page.
- Listed actual enabled Soar, CottenDNS and MasterDNS profile links per delegated DNS zone, with separate copy buttons.
- Kept subscription URL copying separate from direct config copying; no changes to existing account tokens, subscription payloads or quotas.

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
