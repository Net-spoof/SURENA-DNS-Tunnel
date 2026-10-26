<div align="center">
<h1>DNSFabric · Subscriber Portal</h1>
<p><strong>لینک کوتاه، داشبورد اشتراک، QR و نمایش حجم و انقضا</strong></p>
<p><a href="../README.md">صفحه اصلی</a> · <a href="INSTALLATION.fa.md">راهنمای نصب</a> · <a href="SUBSCRIPTION-PORTAL.en.md">English</a></p>
</div>

---

<div dir="rtl" align="right">

<h2>چه چیزی اضافه شده است؟</h2>

<p>این افزونه بر پایه کد رسمی <strong>WhiteDNS Server Management 0.8.2</strong> ساخته می‌شود. برای هر کاربر همان توکن تصادفی امن فعلی را نگه می‌دارد و فقط راه‌های جدیدی برای اشتراک‌گذاری ایجاد می‌کند.</p>

<table>
<thead><tr><th align="right">آدرس</th><th align="right">کاربرد</th><th align="right">سازگاری</th></tr></thead>
<tbody>
<tr><td align="right"><code>https://panel.example.com/sub/TOKEN</code></td><td align="right">ساب رسمی قدیمی با خروجی Base64</td><td align="right">بدون تغییر</td></tr>
<tr><td align="right"><code>https://panel.example.com/s/TOKEN</code></td><td align="right">ساب هوشمند جدید: در مرورگر صفحه کاربر؛ در اپ کانفیگ خام</td><td align="right">همان محتوا و هدرهای ساب قبلی برای اپ</td></tr>
<tr><td align="right"><code>https://panel.example.com/u/TOKEN</code></td><td align="right">صفحه کاربر شامل وضعیت، حجم، انقضا و QR</td><td align="right">برای مرورگر؛ نه جایگزین آدرس ساب در اپ</td></tr>
</tbody>
</table>

<p><strong>تفاوت مهم:</strong> مرورگر با هدرهای مخصوص مرور وب، صفحه HTML را در <code>/s/TOKEN</code> می‌بیند؛ اپ WhiteDNS همان لینک را با خروجی اصلی Base64/Plain دریافت می‌کند. لینک قدیمی <code>/sub/TOKEN</code> برای سازگاری همیشه خام باقی می‌ماند. اگر لینک قدیمی ذخیره کرده‌اید، لینک جدید را از پنل دوباره کپی کنید.</p>
<p><code>TOKEN</code> در مثال بالا فقط جای‌نگهدار است؛ مقدار واقعی آن را پنل تولید می‌کند. این لینک محرمانه است و نباید در کانال عمومی یا لاگ منتشر شود.</p>

<h2>راهنمای بروزرسانی روی سرور</h2>

<p><strong>محدوده پشتیبانی:</strong> نسخه اصلی WhiteDNS روی سرور باید دقیقاً <code>0.8.2</code> باشد. فقط Ubuntu/Debian Linux با معماری amd64 یا arm64 و سرویس systemd استاندارد <code>soar-panel</code> پشتیبانی می‌شود. ابزار به‌صورت fail-closed اجرا می‌شود و روی نسخه‌های دیگر از ادامه کار خودداری می‌کند.</p>

<p>ابتدا اسکریپت را از مخزن خود پروژه دریافت کنید:</p>

</div>

~~~bash
curl -fL --show-error --retry 3 \
  -o /root/dnsfabric-update-subportal.sh \
  https://raw.githubusercontent.com/Net-spoof/dnsfabric/main/update-subportal.sh

bash -n /root/dnsfabric-update-subportal.sh
bash /root/dnsfabric-update-subportal.sh --check
~~~

<div dir="rtl" align="right">

<p>پس از تأیید بررسی سازگاری و وجود نسخه منتشرشده، نصب افزونه:</p>

</div>

~~~bash
bash /root/dnsfabric-update-subportal.sh
~~~

<div dir="rtl" align="right">

<h3>در طول آپدیت چه اتفاقی می‌افتد؟</h3>

<ol>
<li>نسخه باینری موجود، دیتابیس، سرویس و معماری بررسی می‌شوند.</li>
<li>بسته مخصوص معماری سرور از GitHub Release دریافت و <code>SHA256</code> آن بررسی می‌شود.</li>
<li>از باینری و تنظیمات امنیتی بکاپ تهیه می‌شود.</li>
<li>سرویس پنل موقتاً متوقف و بکاپ دیتابیس در مسیر اختصاصی محافظت‌شده ذخیره می‌شود.</li>
<li>فقط باینری پنل عوض و سرویس راه‌اندازی می‌شود؛ فایل‌های کانفیگ تونل و حساب‌ها دستکاری نمی‌شوند.</li>
<li>اگر راه‌اندازی تأیید نشود، باینری قبلی به‌صورت خودکار بازگردانده می‌شود.</li>
</ol>

<p>مسیر بکاپ مشابه <code>/root/dnsfabric-subportal-backup-YYYYMMDDTHHMMSSZ</code> است و به‌صورت پیش‌فرض فقط root اجازه خواندن آن را دارد. <strong>این عملیات یک وقفه کوتاه برای سرویس پنل ایجاد می‌کند</strong>؛ زمان کم‌ترافیک مناسب‌تر است.</p>

<h2>استفاده از پنل</h2>

<ol>
<li>در پنل مدیریت، قسمت <strong>Accounts</strong> را باز کنید.</li>
<li>برای کاربر، روی <strong>لینک / QR</strong> و سپس <strong>لینک ساب</strong> بزنید.</li>
<li>لینک جدید <code>/s/TOKEN</code> برای افزودن به اپ و QR آماده است؛ باز کردن همان لینک در مرورگر به‌صورت خودکار صفحه اختصاصی کاربر را نمایش می‌دهد.</li>
<li>دکمه <strong>نمایش صفحه اشتراک · Subscription portal</strong> صفحه <code>/u/TOKEN</code> را باز می‌کند.</li>
<li>تعویض توکن از همان کنترل قدیمی «تعویض لینک ساب» انجام می‌شود؛ با تعویض، لینک‌های قدیمی و صفحه مربوط به آن‌ها دیگر معتبر نیستند.</li>
</ol>

<h2>حفظ سرویس و امنیت</h2>

<ul>
<li>نیازی به ساخت دوباره کاربران، تعویض کلیدها یا صفرکردن مصرف نیست.</li>
<li>فرمت دانلود <code>plain</code> و <code>base64</code> و هدرهای مصرف و انقضا در آدرس قدیمی بدون تغییر می‌مانند.</li>
<li>صفحه کاربر از توکن اصلی ۳۲بایتی با نمایش URL-safe استفاده می‌کند و مسیر ورود خصوصی ادمین را افشا نمی‌کند.</li>
<li>توکن برای ورود به صفحه عمومی کاربر یک <strong>Bearer Credential</strong> است؛ هر کسی که لینک را داشته باشد می‌تواند وضعیت اشتراک و پروفایل فعال را ببیند.</li>
<li>هدرهای عدم کش و عدم ایندکس در صفحه کاربر اعمال شده‌اند.</li>
<li>آپدیت‌های رسمی WhiteDNS ممکن است این باینری سفارشی را جایگزین کنند. قبل از هر آپدیت اصلی، سازگاری افزونه را مجدداً بررسی کنید.</li>
</ul>

<blockquote>DNSFabric هیچ تغییری در مخزن رسمی WhiteDNS یا پنل OpenVPN/SURENA ایجاد نمی‌کند. این افزونه با حفظ حقوق و مجوز Apache-2.0 بالادستی به‌صورت مستقل ساخته می‌شود؛ انتشار بسته در GitHub به‌معنای تأیید نصب واقعی روی VPS شما نیست.</blockquote>

</div>
