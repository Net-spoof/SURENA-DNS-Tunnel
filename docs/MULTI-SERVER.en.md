<div align="center">
  <h1>DNSFabric · Multiple VPS Deployments</h1>
  <p><strong>Independent DNS delegation, consistent names, repeatable setup</strong></p>
  <p><a href="../README.md">Overview</a> · <a href="INSTALLATION.en.md">Installation</a> · <a href="MULTI-SERVER.fa.md">فارسی</a></p>
</div>

---

## Deployment approach

Use a **unique tunnel subdomain**, **unique authoritative nameserver host**, and **unique HTTPS management hostname** for each VPS. The examples below use the reserved documentation domain <code>example.com</code> and placeholder IP values.

| VPS | Delegated tunnel domain | NS hostname | Panel domain |
|:--|:--|:--|:--|
| 1 | <code>t01.example.com</code> | <code>ns-t01.example.com</code> | <code>panel01.example.com</code> |
| 2 | <code>t02.example.com</code> | <code>ns-t02.example.com</code> | <code>panel02.example.com</code> |
| 3 | <code>t03.example.com</code> | <code>ns-t03.example.com</code> | <code>panel03.example.com</code> |

These names are only examples: location-oriented names such as <code>de1</code>, <code>tr1</code> or <code>eu2</code> are valid if unused and within the **24-character complete tunnel FQDN limit**.

> In this example, <code>example.com</code> is 11 characters long; single-label tunnel prefixes can contain at most 12 characters. The panel and nameserver hostnames are not subject to the tunnel FQDN's 24-character cap.

## Per-server Cloudflare records

Create the following records in the **parent DNS zone**. Replace each IP placeholder with the corresponding VPS's actual public IPv4. Do not overwrite existing records.

### VPS #1

| Type | Name | Content / Target | Proxy |
|:--|:--|:--|:--|
| A | <code>ns-t01</code> | <code>SERVER_1_IP</code> | DNS only |
| NS | <code>t01</code> | <code>ns-t01.example.com</code> | — |
| A | <code>panel01</code> | <code>SERVER_1_IP</code> | DNS only |

### VPS #2

| Type | Name | Content / Target | Proxy |
|:--|:--|:--|:--|
| A | <code>ns-t02</code> | <code>SERVER_2_IP</code> | DNS only |
| NS | <code>t02</code> | <code>ns-t02.example.com</code> | — |
| A | <code>panel02</code> | <code>SERVER_2_IP</code> | DNS only |

### VPS #3

| Type | Name | Content / Target | Proxy |
|:--|:--|:--|:--|
| A | <code>ns-t03</code> | <code>SERVER_3_IP</code> | DNS only |
| NS | <code>t03</code> | <code>ns-t03.example.com</code> | — |
| A | <code>panel03</code> | <code>SERVER_3_IP</code> | DNS only |

The A record for each NS host must resolve to that same server. The NS target is an **FQDN, not an IP address**. Leave the parent's registrar nameservers unchanged.

## Installer answers per VPS

| Prompt | VPS #2 | VPS #3 |
|:--|:--|:--|
| Main domain | <code>example.com</code> | <code>example.com</code> |
| Delegated DNS tunnel domain | <code>t02.example.com</code> | <code>t03.example.com</code> |
| NS hostname | <code>ns-t02.example.com</code> | <code>ns-t03.example.com</code> |
| Panel HTTPS domain | <code>panel02.example.com</code> | <code>panel03.example.com</code> |
| Local IPv4 | Assigned VPS #2 interface IPv4 | Assigned VPS #3 interface IPv4 |
| Public IPv4 | <code>SERVER_2_IP</code> | <code>SERVER_3_IP</code> |
| Administrator password | Dedicated strong secret | Different strong secret |

For each **fresh** VPS, execute the same installer with **that server's** DNS records and IPs:

~~~bash
curl -fL --show-error \
  -o /root/dnsfabric-install.sh \
  https://raw.githubusercontent.com/Net-spoof/dnsfabric/main/install.sh

test -s /root/dnsfabric-install.sh \
  && bash -n /root/dnsfabric-install.sh \
  && bash /root/dnsfabric-install.sh
~~~

### Verify VPS #2 DNS

~~~bash
dig +short A ns-t02.example.com @1.1.1.1
dig +short NS t02.example.com @1.1.1.1
dig +short A panel02.example.com @1.1.1.1
~~~

Both A records must return VPS #2's public IPv4 and NS should return <code>ns-t02.example.com.</code>.

## Operational boundaries

- Each server has **independent accounts, credentials, usage and quotas**. DNSFabric does not sync these across VPS instances.
- DNSFabric does **not** implement automatic cross-server failover or aggregate multiple client connections.
- Do not point a live management domain to a different server without first planning migration and certificates.
- Multiple tunnels are not proof of availability under nationwide internet or international DNS restrictions. Confirm connectivity from each target network.

For first-time configuration, use the [full installation guide](INSTALLATION.en.md).
