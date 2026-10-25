<div dir="rtl" align="right">

# مشارکت در DNSFabric

از مشارکت در بهبود نصب‌کننده، آزمون‌ها، پایداری، مستندات و عیب‌یابی استقبال می‌شود. هدف پروژه، ابزار استقرار مستقل و قابل‌آزمون است، نه انتشار نسخه تغییریافته از WhiteDNS یا موتورهای تحت مجوز دیگر.

## اصول توسعه

<ol>
<li>کد نصب باید <strong>قابل تکرار</strong> باشد و روی سرور فعال سرویس دیگری را بدون تأیید تغییر ندهد.</li>
<li>اعتبارسنجی‌ها و شکست‌های مرتبط با DNS و امنیت باید به‌طور پیش‌فرض <strong>Fail Closed</strong> باشند.</li>
<li>هیچ رمز، Token یا داده شخصی را در کد، Issue، لاگ یا آزمون‌ها وارد نکنید.</li>
<li>در صورت تغییر مسیر نصب، آزمون Mock و حالت <code>--dry-run</code> را بررسی کنید.</li>
<li>مستندات فارسی با راست‌چین استاندارد <code>dir="rtl"</code> و دستورهای ترمینال به‌صورت مستقل چپ‌به‌راست ارائه شوند.</li>
</ol>

## اجرای تست‌ها

</div>

~~~bash
bash -n install.sh
python3 -m py_compile installer.py
python3 -m unittest discover -s tests -v
~~~

<div dir="rtl" align="right">

## ارسال تغییر

۱. یک Branch اختصاصی ایجاد کنید. ۲. تغییر را با تست مرتبط اعمال کنید. ۳. روی سرور آزمایشی تازه بررسی کنید. ۴. Pull Request با شرح محدودیت‌ها، سناریوی تست و اثر روی نصب‌های قبلی ارسال کنید.

> تست Mock اثبات کامل سازگاری با انتشارهای بالادستی یا شبکه‌های محدود نیست. برای هر تغییر عمده در فرایند نصب، آزمایش واقعی و گزارش نسخه محیط لازم است.

</div>

---

## Contributing (English)

Keep installer changes scoped, testable, and fail-closed. Do not commit credentials or upstream proprietary binaries. Run syntax checks, unit tests, and a fresh-VPS integration trial when deployment behavior changes. Describe version compatibility and operational risk in pull requests.
