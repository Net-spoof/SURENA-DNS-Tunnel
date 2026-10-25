<div align="center">
  <h1>نصب DNSFabric</h1>
  <p><strong>راهنمای گام‌به‌گام آماده‌سازی VPS، DNS و پنل WhiteDNS</strong></p>
  <p><a href="../README.md">صفحه اصلی</a> · <a href="INSTALLATION.en.md">English</a> · <a href="MULTI-SERVER.fa.md">چند سرور</a> · <a href="TROUBLESHOOTING.fa.md">رفع خطا</a></p>
</div>

---

<div dir="rtl" align="right">
  <h2>۱. پیش‌نیازها</h2>
  <p>این ابزار برای <strong>سرور تازه و خالی</strong> است، نه به‌روزرسانی یا انتقال پنل موجود.</p>
  <ul>
    <li>Ubuntu 22.04، 24.04، 26.04 یا Debian 12/13 روی معماری amd64 یا arm64.</li>
    <li>SSH تعاملی با دسترسی root و ابزارهای <code>python3</code>، <code>curl</code>، <code>ip</code>، <code>ss</code> و systemd.</li>
    <li>دامنه‌ای که در آن قادر به ساخت رکورد A و NS هستید.</li>
    <li>بازبودن ورودی UDP/53، TCP/53، TCP/80 و TCP/443 در فایروال سیستم‌عامل و شرکت VPS؛ پورت SSH را حفظ کنید.</li>
    <li>اینترنت خارجی روی VPS برای دریافت بسته رسمی WhiteDNS و صدور گواهی TLS.</li>
  </ul>

  <h2>۲. چهار ورودی دامنه چه تفاوتی دارند؟</h2>
  <p>در تمام مثال‌های این راهنما از <strong>دامنه عمومی و رزروشده</strong> <code>example.com</code> استفاده شده است. آدرس‌های مثال، متعلق به کاربر یا یک استقرار واقعی نیستند.</p>
  <table>
    <thead><tr><th align="right">ورودی نصب‌کننده</th><th align="right">معنی</th><th align="right">مثال</th></tr></thead>
    <tbody>
      <tr><td align="right">Main domain</td><td align="right">دامنه اصلی در مدیریت DNS</td><td align="right"><code>example.com</code></td></tr>
      <tr><td align="right">Delegated tunnel domain</td><td align="right">زیردامنه واگذارشده به تونل</td><td align="right"><code>t01.example.com</code></td></tr>
      <tr><td align="right">NS hostname</td><td align="right">نام‌سروری که رکورد A دارد</td><td align="right"><code>ns-t01.example.com</code></td></tr>
      <tr><td align="right">Panel HTTPS domain</td><td align="right">آدرس مستقل ورود پنل</td><td align="right"><code>panel01.example.com</code></td></tr>
    </tbody>
  </table>

  <blockquote><strong>محدودیت نام تونل:</strong> مجموع کاراکترهای نام کامل Tunnel Domain باید حداکثر <strong>۲۴</strong> باشد. برای <code>example.com</code> (۱۱ کاراکتر) حداکثر ۱۲ کاراکتر قبل از نقطه مجاز است. این محدودیت درباره نام پنل یا NS hostname نیست.</blockquote>

  <table>
    <thead><tr><th align="right">نام تونل</th><th align="right">تعداد کاراکتر</th><th align="right">نتیجه</th></tr></thead>
    <tbody>
      <tr><td align="right"><code>t01.example.com</code></td><td align="right">۱۵</td><td align="right">مجاز</td></tr>
      <tr><td align="right"><code>longname1234.example.com</code></td><td align="right">۲۴</td><td align="right">مجاز</td></tr>
      <tr><td align="right"><code>longname12345.example.com</code></td><td align="right">۲۵</td><td align="right">غیرمجاز</td></tr>
    </tbody>
  </table>

  <h2>۳. ساخت رکوردهای Cloudflare</h2>
  <p>به بخش <strong>DNS → Records</strong> دامنه اصلی بروید و سه رکورد نمونه زیر را بسازید. عبارت <code>YOUR_PUBLIC_IP</code> را با IP عمومی <strong>سرور جدید خودتان</strong> جایگزین کنید.</p>
  <table>
    <thead><tr><th align="right">Type</th><th align="right">Name</th><th align="right">Content / Target</th><th align="right">Proxy</th></tr></thead>
    <tbody>
      <tr><td align="right"><code>A</code></td><td align="right"><code>ns-t01</code></td><td align="right"><code>YOUR_PUBLIC_IP</code></td><td align="right">DNS only</td></tr>
      <tr><td align="right"><code>NS</code></td><td align="right"><code>t01</code></td><td align="right"><code>ns-t01.example.com</code></td><td align="right">—</td></tr>
      <tr><td align="right"><code>A</code></td><td align="right"><code>panel01</code></td><td align="right"><code>YOUR_PUBLIC_IP</code></td><td align="right">DNS only</td></tr>
    </tbody>
  </table>
  <p>رکورد <code>NS</code> باید <strong>نام میزبان</strong> را نشان دهد، نه IP. رکورد A نام‌سرور باید به VPS برسد و <strong>DNS only</strong> باشد. رکورد NS دامنه اصلی نزد Registrar را تغییر ندهید. رکوردهای سرویس‌های قبلی را نیز دست‌نخورده نگه دارید.</p>

  <h2>۴. پیدا کردن IPv4 محلی و عمومی</h2>
</div>

~~~bash
ip -4 -br addr
curl -4fsS https://api.ipify.org
echo
~~~

<div dir="rtl" align="right">
  <p>Local IPv4 باید واقعاً روی کارت شبکه VPS تنظیم شده باشد؛ Public IPv4 را با آدرس ثبت‌شده در پنل شرکت سرور مطابقت بدهید. این دو آدرس روی VPS بدون NAT ممکن است یکی باشند.</p>

  <h2>۵. دریافت نصب‌کننده از GitHub</h2>
  <p>دستورها را با کاربر root در یک اتصال SSH تعاملی اجرا کنید؛ فقط متن خام داخل کادر را کپی کنید.</p>
</div>

~~~bash
curl -fL --show-error --retry 3 \
  -o /root/dnsfabric-install.sh \
  https://raw.githubusercontent.com/Net-spoof/dnsfabric/main/install.sh

test -s /root/dnsfabric-install.sh \
  && bash -n /root/dnsfabric-install.sh \
  && bash /root/dnsfabric-install.sh
~~~

<div dir="rtl" align="right">
  <p>آزمون فرم و بررسی مقدماتی بدون نصب سرویس:</p>
</div>

~~~bash
bash /root/dnsfabric-install.sh --dry-run
~~~

<div dir="rtl" align="right">
  <h2>۶. راهنمای پاسخ به سؤال‌های نصب</h2>
  <table>
    <thead><tr><th align="right">سؤال</th><th align="right">پاسخ نمونه</th><th align="right">نکته</th></tr></thead>
    <tbody>
      <tr><td align="right"><code>Main domain you control</code></td><td align="right"><code>example.com</code></td><td align="right">دامنه اصلی، نه آدرس پنل</td></tr>
      <tr><td align="right"><code>Delegated DNS tunnel domain</code></td><td align="right"><code>t01.example.com</code></td><td align="right">همان نامی که رکورد NS برای آن ساخته‌اید</td></tr>
      <tr><td align="right"><code>NS hostname</code></td><td align="right"><code>ns-t01.example.com</code></td><td align="right">مقصد NS با رکورد A مستقل</td></tr>
      <tr><td align="right"><code>Local IPv4 bound on VPS</code></td><td align="right">IP محلی کارت شبکه</td><td align="right">در خروجی دستور <code>ip</code> وجود داشته باشد</td></tr>
      <tr><td align="right"><code>Public IPv4 of VPS</code></td><td align="right">IP عمومی VPS</td><td align="right">مقدار رکوردهای A</td></tr>
      <tr><td align="right"><code>Management panel HTTPS domain</code></td><td align="right"><code>panel01.example.com</code></td><td align="right">آدرس HTTPS مدیریت</td></tr>
      <tr><td align="right"><code>Email for Let's Encrypt</code></td><td align="right">ایمیل معتبر یا Enter</td><td align="right">اختیاری</td></tr>
      <tr><td align="right"><code>Administrator username</code></td><td align="right"><code>admin</code></td><td align="right">نام دلخواه مدیر</td></tr>
      <tr><td align="right"><code>Administrator password</code></td><td align="right">رمز قوی یا Enter</td><td align="right">Enter رمز تصادفی تولید می‌کند</td></tr>
      <tr><td align="right"><code>Type INSTALL</code></td><td align="right"><code>INSTALL</code></td><td align="right">پس از ساخت رکوردها و تنظیم فایروال</td></tr>
    </tbody>
  </table>

  <h2>۷. بررسی رکوردهای منتشرشده</h2>
  <p>دستورهای زیر را پس از ساخت رکوردها اجرا کنید. اگر <code>dig</code> نصب نبود، ابتدا بسته <code>dnsutils</code> را نصب کنید.</p>
</div>

~~~bash
dig +short A ns-t01.example.com @1.1.1.1
dig +short NS t01.example.com @1.1.1.1
dig +short A panel01.example.com @1.1.1.1
~~~

<div dir="rtl" align="right">
  <p>دو خروجی A باید IP عمومی همین VPS را نشان دهند و خروجی NS باید به <code>ns-t01.example.com.</code> اشاره کند. تا تأیید انتشار رکوردها، به مرحله نصب سرویس ادامه ندهید.</p>

  <h2>۸. خروجی پایان نصب و اتصال گوشی</h2>
  <p>پس از تأیید تست‌های نصب رسمی، سالم بودن سرویس‌ها و TLS، اسکریپت پیام <strong>INSTALL SUCCESS</strong>، آدرس کامل پنل با مسیر ورود خصوصی و نام کاربری و رمز مدیر را چاپ می‌کند. آن‌ها را در محل امن نگه دارید و در Issue عمومی قرار ندهید.</p>
  <p>سپس در پنل WhiteDNS، حساب بسازید، یک یا چند موتور را فعال کنید، از بخش Link/QR لینک اشتراک بگیرید و در <a href="https://github.com/WhiteDNS/WhiteDNS-Android/releases">اپ رسمی WhiteDNS Android</a> وارد کنید. هر موتور یک پروفایل قابل انتخاب جداست.</p>

  <h2>۹. نصب روی سرور دوم و سوم</h2>
  <p>برای نام‌گذاری استاندارد، رکوردهای مستقل و جدول ورودی هر VPS، <a href="MULTI-SERVER.fa.md"><strong>راهنمای چندسروری</strong></a> را ببینید.</p>
  <blockquote>تست موفق سمت سرور به معنی تضمین کارکرد در خاموشی اینترنت نیست. اگر Resolver کاربر نتواند درخواست و پاسخ DNS را با Nameserver خارجی مبادله کند، تونل نیز ممکن است قطع شود.</blockquote>
</div>
