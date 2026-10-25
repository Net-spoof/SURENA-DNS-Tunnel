<div align="center">
  <h1>DNSFabric · مدیریت چند VPS</h1>
  <p><strong>نام‌گذاری یکپارچه، رکوردهای مستقل DNS و راه‌اندازی مرحله‌ای</strong></p>
  <p><a href="../README.md">صفحه اصلی</a> · <a href="INSTALLATION.fa.md">نصب</a> · <a href="MULTI-SERVER.en.md">English</a></p>
</div>

---

<div dir="rtl" align="right">
  <h2>ساختار پیشنهادی</h2>
  <p>برای هر سرور یک <strong>دامنه تونل مستقل</strong>، یک <strong>Nameserver</strong> و یک <strong>آدرس پنل HTTPS</strong> تعریف کنید. در این راهنما تمام آدرس‌ها نمونه‌اند و از دامنه رزروشده <code>example.com</code> استفاده می‌کنند.</p>

  <table>
    <thead><tr><th align="right">VPS</th><th align="right">Tunnel Domain</th><th align="right">Nameserver</th><th align="right">Panel Domain</th></tr></thead>
    <tbody>
      <tr><td align="right">سرور اول</td><td align="right"><code>t01.example.com</code></td><td align="right"><code>ns-t01.example.com</code></td><td align="right"><code>panel01.example.com</code></td></tr>
      <tr><td align="right">سرور دوم</td><td align="right"><code>t02.example.com</code></td><td align="right"><code>ns-t02.example.com</code></td><td align="right"><code>panel02.example.com</code></td></tr>
      <tr><td align="right">سرور سوم</td><td align="right"><code>t03.example.com</code></td><td align="right"><code>ns-t03.example.com</code></td><td align="right"><code>panel03.example.com</code></td></tr>
    </tbody>
  </table>

  <p>نام‌ها اختیاری هستند؛ می‌توانید از <code>de1</code>، <code>tr1</code> یا <code>eu2</code> نیز استفاده کنید. تنها شرط این است که رکورد دیگری با همان نام فعال نباشد و طول <strong>دامنه کامل تونل</strong> حداکثر ۲۴ کاراکتر باشد. با <code>example.com</code> که ۱۱ کاراکتر دارد، نام یک‌بخشی تونل حداکثر ۱۲ کاراکتر خواهد بود.</p>

  <h2>رکوردهای Cloudflare برای هر VPS</h2>
  <p>در مثال‌ها <code>SERVER_1_IP</code>، <code>SERVER_2_IP</code> و <code>SERVER_3_IP</code> مقدارهای فرضی‌اند و باید با IP عمومی هر سرور جایگزین شوند.</p>

  <h3>سرور اول</h3>
  <table>
    <thead><tr><th align="right">Type</th><th align="right">Name</th><th align="right">Content</th><th align="right">Proxy</th></tr></thead>
    <tbody>
      <tr><td align="right">A</td><td align="right"><code>ns-t01</code></td><td align="right"><code>SERVER_1_IP</code></td><td align="right">DNS only</td></tr>
      <tr><td align="right">NS</td><td align="right"><code>t01</code></td><td align="right"><code>ns-t01.example.com</code></td><td align="right">—</td></tr>
      <tr><td align="right">A</td><td align="right"><code>panel01</code></td><td align="right"><code>SERVER_1_IP</code></td><td align="right">DNS only</td></tr>
    </tbody>
  </table>

  <h3>سرور دوم</h3>
  <table>
    <thead><tr><th align="right">Type</th><th align="right">Name</th><th align="right">Content</th><th align="right">Proxy</th></tr></thead>
    <tbody>
      <tr><td align="right">A</td><td align="right"><code>ns-t02</code></td><td align="right"><code>SERVER_2_IP</code></td><td align="right">DNS only</td></tr>
      <tr><td align="right">NS</td><td align="right"><code>t02</code></td><td align="right"><code>ns-t02.example.com</code></td><td align="right">—</td></tr>
      <tr><td align="right">A</td><td align="right"><code>panel02</code></td><td align="right"><code>SERVER_2_IP</code></td><td align="right">DNS only</td></tr>
    </tbody>
  </table>

  <h3>سرور سوم</h3>
  <table>
    <thead><tr><th align="right">Type</th><th align="right">Name</th><th align="right">Content</th><th align="right">Proxy</th></tr></thead>
    <tbody>
      <tr><td align="right">A</td><td align="right"><code>ns-t03</code></td><td align="right"><code>SERVER_3_IP</code></td><td align="right">DNS only</td></tr>
      <tr><td align="right">NS</td><td align="right"><code>t03</code></td><td align="right"><code>ns-t03.example.com</code></td><td align="right">—</td></tr>
      <tr><td align="right">A</td><td align="right"><code>panel03</code></td><td align="right"><code>SERVER_3_IP</code></td><td align="right">DNS only</td></tr>
    </tbody>
  </table>

  <blockquote><strong>مهم:</strong> دو سرور نباید به اشتباه از رکورد A با نام یکسان و IP متفاوت استفاده کنند. رکوردهای سرورهای فعال را تغییر ندهید. اگر یک دامنه مدیریت قبلاً اشغال شده، نام پنل دیگری انتخاب کنید.</blockquote>

  <h2>ورودی‌های نصب هر سرور</h2>
  <table>
    <thead><tr><th align="right">فیلد</th><th align="right">VPS دوم</th><th align="right">VPS سوم</th></tr></thead>
    <tbody>
      <tr><td align="right">Main domain</td><td align="right"><code>example.com</code></td><td align="right"><code>example.com</code></td></tr>
      <tr><td align="right">Tunnel domain</td><td align="right"><code>t02.example.com</code></td><td align="right"><code>t03.example.com</code></td></tr>
      <tr><td align="right">NS hostname</td><td align="right"><code>ns-t02.example.com</code></td><td align="right"><code>ns-t03.example.com</code></td></tr>
      <tr><td align="right">Panel domain</td><td align="right"><code>panel02.example.com</code></td><td align="right"><code>panel03.example.com</code></td></tr>
      <tr><td align="right">Local IPv4</td><td align="right">IP محلی سرور دوم</td><td align="right">IP محلی سرور سوم</td></tr>
      <tr><td align="right">Public IPv4</td><td align="right"><code>SERVER_2_IP</code></td><td align="right"><code>SERVER_3_IP</code></td></tr>
    </tbody>
  </table>

  <h2>نصب روی سرور بعدی</h2>
  <p>روی VPS خالی بعدی، با IP و رکوردهای مخصوص همان سرور، از نصب‌کننده مشترک استفاده کنید:</p>
</div>

~~~bash
curl -fL --show-error \
  -o /root/dnsfabric-install.sh \
  https://raw.githubusercontent.com/Net-spoof/dnsfabric/main/install.sh
bash /root/dnsfabric-install.sh
~~~

<div dir="rtl" align="right">
  <h2>اعتبارسنجی رکوردهای VPS دوم</h2>
</div>

~~~bash
dig +short A ns-t02.example.com @1.1.1.1
dig +short NS t02.example.com @1.1.1.1
dig +short A panel02.example.com @1.1.1.1
~~~

<div dir="rtl" align="right">
  <p>رکوردهای A باید IP عمومی سرور دوم را نمایش دهند؛ NS باید به <code>ns-t02.example.com.</code> اشاره کند.</p>
  <blockquote>نصب‌های مستقل در نسخه فعلی حساب‌ها و حجم کاربران را بین VPSها همگام نمی‌کنند. داشتن چند VPS همچنین به‌تنهایی تضمین دسترسی هنگام محدودیت شدید DNS نیست.</blockquote>
</div>
