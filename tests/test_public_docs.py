"""Prevent accidental publication of deployment-specific DNS hostnames in docs."""
from pathlib import Path
import ipaddress
import re
import unittest

ROOT = Path(__file__).resolve().parents[1]
REAL_IR_DOMAIN = re.compile(r"(?i)\b[a-z0-9](?:[a-z0-9.-]*[a-z0-9])?\.ir\b")
IPV4_LITERAL = re.compile(r"(?<![a-z0-9])(?:\d{1,3}\.){3}\d{1,3}(?![a-z0-9])")
DOC_SUFFIXES = {".md", ".svg"}


class PublicDocsHygieneTests(unittest.TestCase):
    def test_public_docs_contain_no_ir_fqdns(self):
        for path in ROOT.rglob("*"):
            if path.is_file() and path.suffix in DOC_SUFFIXES:
                with self.subTest(path=path.relative_to(ROOT)):
                    content = path.read_text(encoding="utf-8")
                    self.assertNotRegex(content, REAL_IR_DOMAIN)

    def test_public_docs_contain_no_unknown_public_ipv4(self):
        allowed_resolvers = {"1.1.1.1", "8.8.8.8", "9.9.9.9"}
        for path in ROOT.rglob("*"):
            if path.is_file() and path.suffix in DOC_SUFFIXES:
                with self.subTest(path=path.relative_to(ROOT)):
                    content = path.read_text(encoding="utf-8")
                    for candidate in IPV4_LITERAL.findall(content):
                        try:
                            ip = ipaddress.ip_address(candidate)
                        except ValueError:
                            continue
                        safe = (
                            candidate in allowed_resolvers
                            or not ip.is_global
                            or ip in ipaddress.ip_network("192.0.2.0/24")
                            or ip in ipaddress.ip_network("198.51.100.0/24")
                            or ip in ipaddress.ip_network("203.0.113.0/24")
                        )
                        self.assertTrue(safe, f"Review example IPv4 in {path}: use a documentation-only address.")


if __name__ == "__main__":
    unittest.main()
