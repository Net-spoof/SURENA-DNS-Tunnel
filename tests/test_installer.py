#!/usr/bin/env python3
"""Offline unit checks, no internet, root or server needed."""
import importlib.util
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location('dnsfabric_installer', ROOT / 'installer.py')
mod = importlib.util.module_from_spec(spec)
spec.loader.exec_module(mod)


class InstallerTests(unittest.TestCase):
    def test_domains(self):
        self.assertEqual(mod.validate_domain('soar.example.com.'), 'soar.example.com')
        for invalid in ('-bad.example.com', 'test..example.com',
                        'https://example.com', 'a' * 64 + '.com'):
            with self.assertRaises(RuntimeError):
                mod.validate_domain(invalid)

    def test_public_ipv4(self):
        self.assertEqual(mod.valid_public_ipv4('9.9.9.9'), '9.9.9.9')
        for invalid in ('192.168.1.1', '10.0.0.1', '127.0.0.1', '203.0.113.1'):
            with self.assertRaises(RuntimeError):
                mod.valid_public_ipv4(invalid)

    def test_expected_prompts(self):
        cfg = dict(zone='soar.example.com', parent='example.com',
                   ns='ns-soar.example.com', local_ip='9.9.9.9',
                   public_ip='9.9.9.9', admin='admin', password='test-password-32')
        answers = mod.answers_for_install(cfg)
        self.assertEqual(len(answers), 15)
        self.assertEqual(answers[0][2], True)
        self.assertEqual(answers[-1][1], 'n')

    def test_pty_wizard(self):
        import tempfile
        import io
        from contextlib import redirect_stdout
        with tempfile.TemporaryDirectory() as temp:
            path = Path(temp) / 'mock-panel-install.sh'
            path.write_text("""#!/bin/bash
set -e
read -r -p 'Tunnel subdomain (example: soar.example.com, max 24 characters): ' zone
read -r -p 'Main DNS zone you can edit [example.com]: ' parent
read -r -p 'Nameserver hostname (inside the main zone, outside the tunnel) [ns-soar.example.com]: ' ns
read -r -p 'Local IPv4 for the DNS listener (behind NAT, use the local IP) [9.9.9.9]: ' lip
read -r -p 'PUBLIC IPv4 from your VPS provider (the A record value) [9.9.9.9]: ' pip
read -r -p 'Resolvers for server-side DNS and tunnel tests [1.1.1.1,8.8.8.8]: ' resolvers
read -r -p 'Public IP or SSH hostname for access instructions [9.9.9.9]: ' ssh_host
read -r -p 'Administrator name [admin]: ' user
read -r -s -p 'New administrator password: ' pass
printf '\\n'
read -r -s -p 'Repeat password: ' pass2
printf '\\n'
read -r -p 'Access mode (1 or 2) [1]: ' access
read -r -p 'Management TCP port [2053]: ' port
read -r -p 'Type INSTALL to install/start this panel, or anything else to cancel: ' confirm
[[ "$confirm" == INSTALL && "$access" == 2 && "$pass" == "$pass2" && "$zone" == soar.example.com ]] || exit 9
echo '[OK] Installation and server-side tunnel checks completed.'
read -r -p 'Open the optional domain/HTTPS setup now? [y/N]: ' https
[[ "$https" == n ]]
""")
            cfg = dict(zone='soar.example.com', parent='example.com',
                       ns='ns-soar.example.com', local_ip='9.9.9.9',
                       public_ip='9.9.9.9', admin='admin', password='test-password-32')
            capture = io.StringIO()
            with redirect_stdout(capture):
                log = mod.run_install_pty(path, cfg)
            self.assertIn('Installation and server-side tunnel checks completed.', log)
            self.assertNotIn(cfg['password'], capture.getvalue())

    def test_pty_https_security(self):
        import tempfile
        import io
        from contextlib import redirect_stdout
        with tempfile.TemporaryDirectory() as temp:
            script = Path(temp) / 'mock-security.sh'
            script.write_text("""#!/bin/bash
set -e
read -r -p 'Choose [1]: ' opt
read -r -p 'Panel domain (leave empty to keep the existing domain, or use only a random path): ' domain
read -r -p 'Certificate account email (optional): ' email
read -r -p 'Generate a new random login path? [Y/n]: ' random
read -r -p 'Type APPLY to continue: ' confirm
[[ "$opt" == 1 && "$domain" == panel.example.com && "$random" == Y && "$confirm" == APPLY ]] || exit 4
echo 'Secure access settings saved. Bookmark this private login address:'
echo 'https://panel.example.com/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/'
""")
            responses = [
                (r'Choose \[1\]:', '1'),
                (r'Panel domain \(leave empty to keep the existing domain, or use only a random path\):', 'panel.example.com'),
                (r'Certificate account email \(optional\):', ''),
                (r'Generate a new random login path\? \[Y/n\]:', 'Y'),
                (r'Type APPLY to continue:', 'APPLY'),
            ]
            with redirect_stdout(io.StringIO()):
                result = mod.PtyScript(['bash', str(script)], responses, timeout=5).run()
            self.assertIn('Secure access settings saved.', result)

    def test_login_url(self):
        url = 'https://panel.example.com/' + 'a' * 43 + '/'
        self.assertEqual(mod.URL_RE.findall('\nURL: ' + url), [url])


if __name__ == '__main__':
    unittest.main()
