#!/usr/bin/env python3
"""DNSFabric 1.0 — guided, fail-closed WhiteDNS panel deployment.

No patching of WhiteDNS, no cloud credentials, no logging of passwords.
Pins and checksum-verifies the tested upstream 0.8.2 release. Only new servers.
"""
import argparse
import getpass
import hashlib
import ipaddress
import os
from pathlib import Path, PurePosixPath
import pty
import re
import secrets
import select
import shutil
import signal
import socket
import ssl
import subprocess
import sys
import tarfile
import tempfile
import time

VERSION = "0.8.2"
REPO = "WhiteDNS/Soar-Server"
MGMT_PORT = 2053
URL_RE = re.compile(r"https://[a-z0-9.-]+/[A-Za-z0-9_-]{15,}/")
ANSI = re.compile(r"\x1b\[[0-9;]*[A-Za-z]")


def abort(message):
    raise RuntimeError(message)


def ask(label, default=None):
    hint = f" [{default}]" if default is not None else ""
    val = input(f"{label}{hint}: ").strip()
    return val or default or ""


def validate_domain(value):
    value = value.rstrip('.').lower()
    if len(value) > 253 or not value or any(
        not re.fullmatch(r"[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?", label)
        for label in value.split('.')
    ):
        abort("Invalid ASCII domain: " + value)
    return value


def valid_public_ipv4(value):
    addr = ipaddress.ip_address(value)
    if addr.version != 4 or not addr.is_global:
        abort("Public IPv4 must be globally routed: " + value)
    return str(addr)


def local_ip():
    r = subprocess.run(["ip", "-4", "route", "get", "1.1.1.1"], capture_output=True, text=True, check=True)
    match = re.search(r"\bsrc\s+(\d+(?:\.\d+){3})\b", r.stdout)
    if not match:
        abort("Could not identify local IPv4 from the default route")
    return match.group(1)


def check_host():
    if os.geteuid() != 0:
        abort("Run as root")
    if not sys.stdin.isatty() or not sys.stdout.isatty():
        abort("Use an interactive SSH session (a TTY is required)")
    props = {}
    for line in Path('/etc/os-release').read_text().splitlines():
        if '=' in line:
            key, val = line.split('=', 1)
            props[key] = val.strip('"')
    if props.get('ID') not in ('ubuntu', 'debian') or (props.get('ID'), props.get('VERSION_ID')) not in {
        ('ubuntu', '22.04'), ('ubuntu', '24.04'), ('ubuntu', '26.04'),
        ('debian', '12'), ('debian', '13'),
    }:
        abort("Supported: Ubuntu 22.04/24.04/26.04, Debian 12/13")
    arch = os.uname().machine
    if arch not in ('x86_64', 'amd64', 'aarch64', 'arm64'):
        abort('Only amd64/arm64 supported')
    arch = 'arm64' if arch in ('aarch64', 'arm64') else 'amd64'
    for cmd in ('curl', 'tar', 'sha256sum', 'systemctl', 'ss', 'ip'):
        if not shutil.which(cmd):
            abort(f'Missing required tool: {cmd}. Install it first.')
    if Path('/var/lib/soar-panel/panel.db').exists() or Path('/etc/systemd/system/soar-panel.service').exists():
        abort('Existing WhiteDNS panel detected. No overwrite or upgrade through this installer.')
    if Path('/etc/systemd/system/soar-server.service').exists():
        abort('Existing standalone Soar detected. Use a separate VPS; no automatic migration.')
    if shutil.disk_usage('/').free < 1_500_000_000:
        abort('Less than 1.5 GB free disk space.')
    return arch


def check_ports(address):
    r = subprocess.run(['ss', '-H', '-lntup'], capture_output=True, text=True, check=True)
    collisions = []
    for line in r.stdout.splitlines():
        fields = line.split()
        if len(fields) < 5:
            continue
        hostport = fields[4]
        for port in (53, 80, 443, MGMT_PORT):
            if not hostport.endswith(':' + str(port)):
                continue
            bind = hostport.rsplit(':', 1)[0].split('%')[0].strip('[]')
            # The host may use local DNS stubs (127.0.0.53:53).
            if port == 53 and bind.startswith('127.') and bind != '127.0.0.1':
                continue
            if port == MGMT_PORT and bind != '127.0.0.1' and bind not in ('0.0.0.0', '*'):
                continue
            if port == 53 and bind not in (address, '0.0.0.0', '*'):
                continue
            collisions.append(line)
    if collisions:
        abort('Conflicting listeners (resolve before setup):\n' + '\n'.join(collisions))


def download(url, destination):
    subprocess.run([
        'curl', '--fail', '--silent', '--show-error', '--location', '--retry', '3',
        '--connect-timeout', '20', '--max-time', '240', '--proto', '=https',
        '--proto-redir', '=https', url, '-o', str(destination)
    ], check=True)


def file_sha256(path):
    hash_ = hashlib.sha256()
    with open(path, 'rb') as src:
        while chunk := src.read(1024 * 1024):
            hash_.update(chunk)
    return hash_.hexdigest()


def get_official_installer(tmp, arch):
    asset = f'soar-panel-{VERSION}-linux-{arch}.tar.gz'
    base = f'https://github.com/{REPO}/releases/download/panel-v{VERSION}'
    tarball, sha = tmp / asset, tmp / (asset + '.sha256')
    print(f'\nDownloading WhiteDNS official {VERSION} for {arch}...')
    download(f'{base}/{asset}', tarball)
    download(f'{base}/{asset}.sha256', sha)
    lines = sha.read_text('ascii').splitlines()
    if len(lines) != 1:
        abort('Invalid SHA-256 manifest')
    match = re.fullmatch(r'([a-f0-9]{64})\s+\*?([^\s]+)', lines[0])
    if not match or match.group(2) != asset:
        abort('Unexpected filename or checksum format')
    if file_sha256(tarball) != match.group(1):
        abort('Checksum mismatch: refusing to install')
    root = f'soar-panel-{VERSION}-linux-{arch}'
    with tarfile.open(tarball, 'r:gz') as tar:
        members = tar.getmembers()
        if len(members) > 2000 or sum(m.size for m in members) > 180_000_000:
            abort('Unexpectedly large upstream archive')
        for member in members:
            path = PurePosixPath(member.name)
            if (not member.isfile() and not member.isdir()) or path.is_absolute() or '..' in path.parts or path.parts[0] != root:
                abort('Unsafe archive content: refusing extraction')
        if sys.version_info >= (3, 12):
            tar.extractall(tmp, filter='data')
        else:
            tar.extractall(tmp)
    installer = tmp / root / 'panel-install.sh'
    subprocess.run(['bash', str(installer), '--check-bundle'], check=True)
    return installer


class PtyScript:
    """Only send exact expected answers to an official, checksum-verified wizard."""
    def __init__(self, command, responses, timeout=900):
        self.command, self.responses, self.timeout = command, responses, timeout
        self.transcript = ''

    def run(self):
        pid, fd = pty.fork()
        if pid == 0:
            os.execvp(self.command[0], self.command)
        tail = ''
        pending = list(self.responses)
        timeout_at = time.monotonic() + self.timeout
        error = None
        try:
            while True:
                if time.monotonic() >= timeout_at:
                    error = 'Installer timeout; inspect status before retrying.'
                    break
                r, _, _ = select.select([fd], [], [], min(1, max(0, timeout_at - time.monotonic())))
                if r:
                    try:
                        chunk = os.read(fd, 32768)
                    except OSError:
                        chunk = b''
                    if not chunk:
                        break
                    text_chunk = chunk.decode('utf-8', 'replace')
                    self.transcript = (self.transcript + text_chunk)[-50000:]
                    sys.stdout.write(text_chunk)
                    sys.stdout.flush()
                    tail = ANSI.sub('', (tail + text_chunk)[-4000:]).replace('\r', '')
                    if pending:
                        pattern, response = pending[0]
                        if re.search(pattern + r'\s*$', tail):
                            os.write(fd, (response + '\n').encode('utf-8'))
                            pending.pop(0)
                            tail = ''
                    if re.search(r'Choose \[1\]:\s*$', tail):
                        error = 'Official DNS/health check requested manual intervention. Aborting without bypassing verification.'
                        os.write(fd, b'4\n')
                        break
                    if 'Type MIGRATE to accept this migration' in tail:
                        error = 'Existing standalone installation detected: migration forbidden.'
                        break
            if error:
                try:
                    os.kill(pid, signal.SIGTERM)
                except ProcessLookupError:
                    pass
            _, status = os.waitpid(pid, 0)
            code = os.waitstatus_to_exitcode(status)
            if error:
                abort(error)
            if code != 0:
                abort(f'Official installer exited with code {code}. Do not assume success.')
            if pending:
                abort(f'Installer interaction changed; {len(pending)} expected prompts were not encountered. Verify installation manually.')
            return self.transcript
        finally:
            os.close(fd)


def answers_for_install(cfg):
    items = [
        (r'Install these packages using apt-get\? \(Y/n\) \[Y\]:', 'Y', True),
        (r'Tunnel subdomain \(example: soar\.example\.com, max 24 characters\):', cfg['zone'], False),
        (r'Main DNS zone you can edit(?: \[[^\]]+\])?:', cfg['parent'], False),
        (r'Nameserver hostname \(inside the main zone, outside the tunnel\)(?: \[[^\]]+\])?:', cfg['ns'], False),
        (r'Local IPv4 for the DNS listener \(behind NAT, use the local IP\)(?: \[[^\]]+\])?:', cfg['local_ip'], False),
        (r'PUBLIC IPv4 from your VPS provider \(the A record value\)(?: \[[^\]]+\])?:', cfg['public_ip'], False),
        (r'Resolvers for server-side DNS and tunnel tests(?: \[[^\]]+\])?:', '1.1.1.1,8.8.8.8', False),
        (r'Public IP or SSH hostname for access instructions(?: \[[^\]]+\])?:', cfg['public_ip'], False),
        (r'Administrator name(?: \[[^\]]+\])?:', cfg['admin'], False),
        (r'New administrator password:', cfg['password'], False),
        (r'Repeat password:', cfg['password'], False),
        (r'Access mode \(1 or 2\)(?: \[[^\]]+\])?:', '2', False),
        (r'Management TCP port(?: \[[^\]]+\])?:', str(MGMT_PORT), False),
        (r'Type INSTALL to install/start this panel, or anything else to cancel:', 'INSTALL', False),
        (r'Open the optional domain/HTTPS setup now\? \[y/N\]:', 'n', False),
    ]
    return items


def run_install_pty(installer, cfg):
    full = answers_for_install(cfg)
    expected = [(pattern, response) for pattern, response, optional in full if not optional]
    return run_with_optional(['bash', str(installer)], expected, full[0][:2])


def run_with_optional(command, responses, optional=None):
    """Drop optional apt prompt when first mandatory prompt arrives."""
    if optional is None:
        return PtyScript(command, responses).run()
    class OptionalPty(PtyScript):
        def run(self):
            pid, fd = pty.fork()
            if pid == 0:
                os.execvp(self.command[0], self.command)
            pending = list(self.responses)
            transcript, tail = '', ''
            optional_handled = False
            deadline = time.monotonic() + self.timeout
            error = None
            try:
                while True:
                    if time.monotonic() >= deadline:
                        error = 'Installer timeout; examine service state.'
                        break
                    ready, _, _ = select.select([fd], [], [], 1)
                    if not ready:
                        continue
                    try:
                        chunk = os.read(fd, 32768)
                    except OSError:
                        chunk = b''
                    if not chunk:
                        break
                    text_chunk = chunk.decode('utf-8', 'replace')
                    transcript = (transcript + text_chunk)[-75000:]
                    sys.stdout.write(text_chunk)
                    sys.stdout.flush()
                    tail = ANSI.sub('', (tail + text_chunk)[-4000:]).replace('\r', '')
                    if not optional_handled and re.search(optional[0] + r'\s*$', tail):
                        os.write(fd, (optional[1] + '\n').encode())
                        optional_handled = True
                        tail = ''
                        continue
                    if pending and re.search(pending[0][0] + r'\s*$', tail):
                        os.write(fd, (pending[0][1] + '\n').encode('utf-8'))
                        pending.pop(0)
                        tail = ''
                        continue
                    if re.search(r'Choose \[1\]:\s*$', tail):
                        error = 'Installer requested manual DNS/test intervention; no unverified installation will be accepted.'
                        os.write(fd, b'4\n')
                        break
                    if 'Type MIGRATE to accept this migration' in tail:
                        error = 'Migration prompt unexpected; refusing to migrate an existing service.'
                        break
                if error:
                    try:
                        os.kill(pid, signal.SIGTERM)
                    except ProcessLookupError:
                        pass
                _, status = os.waitpid(pid, 0)
                if error:
                    abort(error)
                if os.waitstatus_to_exitcode(status) != 0:
                    abort('Official installer returned an error; installation not verified.')
                if pending:
                    abort('Installer prompts differed from pinned version; success not established.')
                if '[OK] Installation and server-side tunnel checks completed.' not in transcript:
                    abort('Official server-side tunnel tests did not report success.')
                return transcript
            finally:
                os.close(fd)
    return OptionalPty(command, responses).run()


def run_security(cfg):
    print('\nConfiguring HTTPS and a new private login path...')
    transcript = PtyScript(['soar-panel-security'], [
        (r'Choose \[1\]:', '1'),
        (r'Panel domain \(leave empty to keep the existing domain, or use only a random path\):', cfg['panel']),
        (r'Certificate account email \(optional\):', cfg['email']),
        (r'Generate a new random login path\? \[Y/n\]:', 'Y'),
        (r'Type APPLY to continue:', 'APPLY'),
    ], timeout=360).run()
    matches = URL_RE.findall(transcript)
    url = next((x for x in reversed(matches) if x.startswith('https://' + cfg['panel'] + '/')), None)
    if not url or 'Secure access settings saved.' not in transcript:
        abort('HTTPS or private access setup did not produce a verified login URL.')
    return url


def verify_result(cfg):
    for service in ('soar-panel', 'soar-panel-ops'):
        r = subprocess.run(['systemctl', 'is-active', '--quiet', service], check=False)
        if r.returncode:
            abort(service + ' is not active')
    ctx = ssl.create_default_context()
    with socket.create_connection((cfg['local_ip'], 443), timeout=10) as sock:
        with ctx.wrap_socket(sock, server_hostname=cfg['panel']) as tls:
            if tls.version() not in ('TLSv1.2', 'TLSv1.3'):
                abort('HTTPS TLS version unexpected')


def collect_config():
    print('\nDNSFabric — WhiteDNS guided setup (FA / EN)')
    print('Installation is for a FRESH VPS only. Existing installations will NOT be replaced.')
    guessed = local_ip()
    parent = validate_domain(ask('Main domain you control (e.g. example.com)'))
    default_zone = 'soar.' + parent
    zone = validate_domain(ask('Delegated DNS tunnel domain', default_zone))
    if len(zone) > 24 or not zone.endswith('.' + parent):
        abort('Tunnel domain must be under the main domain and at most 24 characters')
    ns = validate_domain(ask('NS hostname (not under the tunnel)', 'ns-soar.' + parent))
    if not ns.endswith('.' + parent) or ns == zone or ns.endswith('.' + zone):
        abort('NS host must be under main domain but outside delegated tunnel zone')
    local = str(ipaddress.IPv4Address(ask('Local IPv4 bound on VPS', guessed)))
    assigned = subprocess.run(['ip', '-o', '-4', 'addr', 'show'], capture_output=True, text=True, check=True).stdout
    if not re.search(r'\binet\s+' + re.escape(local) + r'/', assigned):
        abort('Local IPv4 is not assigned to this VPS')
    public = valid_public_ipv4(ask('Public IPv4 of VPS', local))
    panel = validate_domain(ask('Management panel HTTPS domain', 'panel.' + parent))
    if panel == zone or panel.endswith('.' + zone):
        abort('Panel domain must be outside delegated DNS tunnel')
    email = ask("Email for Let's Encrypt (optional)")
    if email and (len(email) > 254 or not re.fullmatch(r'[^\s@]+@[^\s@]+\.[^\s@]+', email)):
        abort('Invalid certificate email')
    admin = ask('Panel administrator username', 'admin')
    if not re.fullmatch(r'[A-Za-z0-9_.-]{1,64}', admin):
        abort('Invalid administrator name')
    p = getpass.getpass('Administrator password (Enter = generate random): ')
    if not p:
        p = secrets.token_urlsafe(24)
        print('A strong random administrator password was generated; it will be displayed at the END.')
    if not 12 <= len(p) <= 72 or not all(32 <= ord(c) <= 126 for c in p):
        abort('Password must be 12–72 printable ASCII characters')
    cfg = {'parent': parent, 'zone': zone, 'ns': ns, 'local_ip': local,
           'public_ip': public, 'panel': panel, 'email': email, 'admin': admin, 'password': p}
    check_ports(local)
    print('\n===== REQUIRED DNS RECORDS (Cloudflare: DNS only) =====')
    print(f'A  {ns}  ->  {public}')
    print(f'NS {zone}  ->  {ns}')
    print(f'A  {panel}  ->  {public}')
    print('Do NOT change the main domain nameservers or existing records.')
    print('Provider firewall: allow UDP/TCP 53; TCP 80 and 443; keep SSH reachable.')
    if ask('After creating the above records, type INSTALL to proceed') != 'INSTALL':
        abort('Cancelled before installation')
    return cfg


def main():
    parser = argparse.ArgumentParser(description='DNSFabric WhiteDNS guided installer')
    parser.add_argument('--dry-run', action='store_true', help='validate prompts and show guidance without installing')
    args = parser.parse_args()
    arch = check_host()
    cfg = collect_config()
    if args.dry_run:
        print('\n[DRY RUN] All input checks passed. No packages/services modified.')
        return
    with tempfile.TemporaryDirectory(prefix='dnsfabric-whitedns-', dir='/tmp') as d:
        installer = get_official_installer(Path(d), arch)
        print('\nSTARTING OFFICIAL INSTALLER. No secrets are written to disk by DNSFabric.')
        run_install_pty(installer, cfg)
    print('\nOfficial panel installed and server-side probes passed.')
    url = run_security(cfg)
    verify_result(cfg)
    print('\n' + '=' * 64)
    print('INSTALL SUCCESS — DNSFabric / WhiteDNS ' + VERSION)
    print('Panel URL : ' + url)
    print('Username  : ' + cfg['admin'])
    print('Password  : ' + cfg['password'])
    print('Server IP : ' + cfg['public_ip'])
    print('Tunnel    : ' + cfg['zone'])
    print('Engines   : Soar, CottenDNS, MasterDNS (available per account)')
    print('Next: log in, create account, import subscription into WhiteDNS Android.')
    print('IMPORTANT: Save password privately. It is intentionally not saved by this wrapper.')
    print('This confirms server-side tests, NOT access during a nationwide shutdown.')
    print('=' * 64)


if __name__ == '__main__':
    try:
        main()
    except (RuntimeError, ValueError, subprocess.CalledProcessError, OSError, EOFError, KeyboardInterrupt) as exc:
        print('\nINSTALL NOT VERIFIED: ' + str(exc), file=sys.stderr)
        print('Inspect systemctl status soar-panel and journalctl -u soar-panel before retrying.', file=sys.stderr)
        sys.exit(1)
