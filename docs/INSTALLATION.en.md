<div align="center">
  <h1>DNSFabric · Installation</h1>
  <p><strong>From a clean VPS to validated WhiteDNS HTTPS management</strong></p>
  <p><a href="../README.md">Overview</a> · <a href="INSTALLATION.fa.md">فارسی</a> · <a href="MULTI-SERVER.en.md">Multiple VPS</a> · <a href="TROUBLESHOOTING.fa.md">Troubleshooting</a></p>
</div>

---

## 1 · Requirements

DNSFabric v1.0.0 targets **fresh** Ubuntu 22.04/24.04/26.04 or Debian 12/13 servers (amd64/arm64). You need root access through interactive SSH, Python 3.10+, curl, systemd, publicly reachable IPv4 and a domain whose DNS records you control.

Allow **UDP/TCP 53** for DNS tunneling and **TCP 80/443** for HTTPS management and certificate issuance. Retain SSH access and check both OS-level and provider firewalls.

> Do not run this fresh-install wrapper against an existing WhiteDNS/Soar deployment. It does not support upgrade or migration.

## 2 · Understand the four hostnames

All examples use the **reserved documentation domain** <code>example.com</code>, not a real deployment.

| Installer input | Purpose | Example |
|:--|:--|:--|
| Main domain | Parent zone managed in Cloudflare | <code>example.com</code> |
| Delegated DNS tunnel domain | Child zone containing the tunnel | <code>t01.example.com</code> |
| NS hostname | Authoritative nameserver host | <code>ns-t01.example.com</code> |
| Management panel HTTPS domain | Management endpoint | <code>panel01.example.com</code> |

**Tunnel domain length limit:** The complete tunnel FQDN must be **24 characters or fewer**. The main domain <code>example.com</code> is 11 characters; a single-label prefix can be at most 12 characters.

| Example | Length | Valid? |
|:--|--:|:--|
| <code>t01.example.com</code> | 15 | Yes |
| <code>longname1234.example.com</code> | 24 | Yes |
| <code>longname12345.example.com</code> | 25 | No |

The 24-character limit applies to the tunnel domain only, **not** the panel hostname or NS hostname.

## 3 · Configure Cloudflare DNS

Open **Cloudflare → Domain → DNS → Records** and add these entries to the parent zone. Replace <code>YOUR_PUBLIC_IP</code> with **your new server's public IPv4**.

| Type | Name | Value / Target | Proxy |
|:--|:--|:--|:--|
| A | <code>ns-t01</code> | <code>YOUR_PUBLIC_IP</code> | DNS only |
| NS | <code>t01</code> | <code>ns-t01.example.com</code> | — |
| A | <code>panel01</code> | <code>YOUR_PUBLIC_IP</code> | DNS only |

The <code>NS</code> record delegates <code>t01.example.com</code> to <code>ns-t01.example.com</code>; the nameserver's A record points to the VPS. The NS record value must be a hostname, **not an IP address**. Keep the A records DNS-only (grey cloud).

Do not change your parent domain's registrar nameservers or accidentally overwrite records belonging to another service.

## 4 · Check local and public addresses

~~~bash
ip -4 -br addr
curl -4fsS https://api.ipify.org
echo
~~~

Use the IPv4 actually assigned to the VPS interface for the **local** listener. Public IPv4 must match the address supplied by the VPS provider and the target of both A records. The local and public addresses may differ behind NAT.

## 5 · Install from the official DNSFabric repository

Run under an interactive SSH terminal as root:

~~~bash
curl -fL --show-error --retry 3 \
  -o /root/dnsfabric-install.sh \
  https://raw.githubusercontent.com/Net-spoof/dnsfabric/main/install.sh

test -s /root/dnsfabric-install.sh \
  && bash -n /root/dnsfabric-install.sh \
  && bash /root/dnsfabric-install.sh
~~~

The wrapper downloads the pinned upstream **WhiteDNS 0.8.2** release and validates its SHA-256 file before executing the official installer.

To check the inputs and initial host requirements without installing services:

~~~bash
bash /root/dnsfabric-install.sh --dry-run
~~~

## 6 · Interactive installer answers

| Prompt | Example answer | Note |
|:--|:--|:--|
| Main domain you control | <code>example.com</code> | The parent zone |
| Delegated DNS tunnel domain | <code>t01.example.com</code> | The delegated child domain |
| NS hostname | <code>ns-t01.example.com</code> | The NS target with an A record |
| Local IPv4 bound on VPS | IPv4 assigned to network interface | Not necessarily the public address |
| Public IPv4 of VPS | Public IPv4 from your VPS provider | Target of A records |
| Management panel HTTPS domain | <code>panel01.example.com</code> | Separate management hostname |
| Email for Let's Encrypt | Valid email or Enter | Optional |
| Administrator username | <code>admin</code> | Your chosen login |
| Administrator password | Strong password or Enter | Enter generates a random password |
| Type INSTALL | <code>INSTALL</code> | Only after DNS / ports are ready |

The wrapper answers pinned upstream prompts. An unexpected prompt or failing test must not be treated as a successful installation.

## 7 · Check DNS propagation

~~~bash
dig +short A ns-t01.example.com @1.1.1.1
dig +short NS t01.example.com @1.1.1.1
dig +short A panel01.example.com @1.1.1.1
~~~

Both A records should resolve to the **new VPS public IPv4**. The NS query should return <code>ns-t01.example.com.</code>. A trailing dot in the DNS response is normal. Wait for propagation or fix conflicting DNS records before proceeding.

## 8 · Success, credentials, and Android

DNSFabric reports **INSTALL SUCCESS** only after the official server-side checks, service health and TLS validation succeed. Save the complete private panel URL and administrator credentials in a secure password manager. Do not post them publicly.

Within the panel, create an account, select Soar/CottenDNS/MasterDNS engines, obtain the Subscription URL through Link/QR and import it into the [official WhiteDNS Android client](https://github.com/WhiteDNS/WhiteDNS-Android/releases). Each engine remains an independently selectable profile.

> DNS tunneling depends on the user's resolver being able to exchange DNS queries and responses with the remote authoritative nameserver. Successful server-side installation does not guarantee functionality during a complete internet shutdown.

## 9 · Scaling to multiple VPS instances

For consistent names, independent records and worked values for a second/third server, see [Multi-server deployment](MULTI-SERVER.en.md).
