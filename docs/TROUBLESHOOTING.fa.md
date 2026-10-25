<div align="center">

<h1>DNSFabric — راهنمای عیب‌یابی</h1>

<p>بررسی مرحله‌ای خطاهای DNS، نصب، HTTPS و اتصال اندروید</p>

<p><a href="../README.md">بازگشت به راهنمای اصلی</a></p>

</div>

---

<div dir="rtl" align="right">

<h2>۱. نصب در بررسی DNS متوقف شده است</h2>

<p>قبل از اجرای مجدد، این موارد را بررسی کنید:</p>
<ul>
<li>رکورد <code>A</code> نیم‌سرور در Cloudflare به IP عمومی درست VPS اشاره کند و روی <strong>DNS only</strong> باشد.</li>
<li>رکورد <code>NS</code> دامنه تونل، نیم‌سرور همان دامنه را معرفی کند.</li>
<li>نام کامل زیردامنه DNS Tunnel حداکثر ۲۴ کاراکتر باشد.</li>
<li>پورت‌های <code>UDP/53</code> و <code>TCP/53</code> در فایروال VPS و ارائه‌دهنده باز باشند.</li>
<li>در صورت NAT، IP محلی برای Listener و IP عمومی برای DNS درست و متفاوت وارد شوند.</li>
</ul>

<blockquote><strong>هشدار:</strong> شکست نصب به معنی پاک‌بودن سرور نیست. قبل از تلاش دوباره، وضعیت فایل‌ها و سرویس‌ها را بررسی کنید. نصب‌کننده DNSFabric روی سرور دارای دیتابیس پنل قبلی متوقف می‌شود.</blockquote>

</div>

~~~bash
systemctl status soar-panel soar-panel-ops --no-pager
journalctl -u soar-panel -n 80 --no-pager
test -e /var/lib/soar-panel/panel.db && echo "Panel database exists"
~~~

<div dir="rtl" align="right">

<h2>۲. موتورهای DNS نصب شده‌اند ولی HTTPS پنل ناموفق است</h2>

<ul>
<li>رکورد <code>A</code> دامنه پنل روی IP عمومی VPS و بدون Cloudflare Proxy تنظیم شود.</li>
<li>رکورد <code>AAAA</code> اشتباه یا تنظیمات DNS قدیمی را بررسی کنید.</li>
<li>ورودی <code>TCP/80</code> و <code>TCP/443</code> باید برای صدور و تمدید گواهی در دسترس باشد.</li>
<li>وجود Nginx/Apache یا سرویس دیگری روی ۸۰/۴۴۳ می‌تواند مانع تنظیم HTTPS شود.</li>
</ul>

<p>برای تکمیل یا بازیابی HTTPS از ابزار رسمی نصب‌شده استفاده کنید:</p>

</div>

~~~bash
sudo soar-panel-security
~~~

<div dir="rtl" align="right">

<h2>۳. آدرس کامل ورود به پنل را ندارم</h2>

<p>از منوی ابزار امنیت، گزینه <strong>۲</strong> را برای نمایش آدرس خصوصی فعلی انتخاب کنید. مسیر تصادفی بعد از دامنه بخشی از آدرس ورود است و باید حفظ شود.</p>

</div>

~~~bash
sudo soar-panel-security
# Select option 2: Show the current login address
~~~

<div dir="rtl" align="right">

<h2>۴. Soar وصل می‌شود ولی CottenDNS یا MasterDNS نه</h2>

<ul>
<li>از پنل مطمئن شوید موتورهای موردنظر در همان حساب فعال‌اند.</li>
<li>پس از ویرایش حساب، در WhiteDNS Android اتصال را قطع و <strong>Refresh subscriptions</strong> را اجرا کنید.</li>
<li>برای هر موتور، پروفایل جداگانه انتخاب کنید؛ اپ به‌صورت خودکار همه موتورهای یک اشتراک را ترکیب نمی‌کند.</li>
<li>DNS Resolverهای همان شبکه موبایل را اسکن کنید.</li>
<li>لاگ اتصال، Timeout و تست واقعی مرورگر و IP خروجی را بررسی کنید.</li>
</ul>

<h2>۵. وضعیت برنامه Connected است ولی اینترنت باز نمی‌شود</h2>

<p>Connected بودن لزوماً به معنی انتقال موفق ترافیک نیست. یک درخواست HTTPS و IP خروجی را بررسی کنید. تفاوت حالت Proxy و Full VPN، تنظیمات Split Tunnel و وضعیت Resolverها اهمیت دارد.</p>

<h2>۶. سرور قبلاً WhiteDNS یا Soar داشته است</h2>

<p>DNSFabric فعلاً فقط نصب اولیه روی VPS تازه را پشتیبانی می‌کند. این ابزار جایگزین عملیات <strong>آپدیت، مهاجرت و بازیابی حساب‌ها</strong> نیست. برای ارتقا از مسیر رسمی پنل و بکاپ استفاده کنید.</p>

<h2>۷. قطعی کامل اینترنت و محدودیت DNS</h2>

<p>اگر Resolver اپراتور درخواست‌ها را به نام‌سرور خارجی Forward نکند یا پاسخ را برنگرداند، هیچ تضمینی برای اتصال DNS Tunnel وجود ندارد. عبور موفق در شرایط معمول، نتیجه آزمایش در قطعی کامل اینترنت نیست.</p>

<h2>اطلاعات مناسب برای گزارش خطا</h2>

<ul>
<li>نسخه سیستم‌عامل و معماری سرور</li>
<li>مرحله‌ای از نصب که متوقف شده، به‌همراه متن خطا</li>
<li>وضعیت <code>soar-panel</code> و <code>soar-panel-ops</code></li>
<li>نتیجه DNS Delegation و فعال‌بودن پورت‌های مربوطه</li>
<li>نوع اتصال کلاینت، موتور انتخاب‌شده و گزارش بدون اطلاعات حساس</li>
</ul>

<p><strong>هرگز</strong> رمز مدیر، توکن Subscription، کلید خصوصی یا URL خصوصی ورود پنل را داخل Issue عمومی منتشر نکنید.</p>

</div>
