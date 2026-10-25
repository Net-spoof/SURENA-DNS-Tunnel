<div dir="rtl" align="right">

# سیاست امنیتی DNSFabric

DNSFabric ابزار استقرار و نصب اولیه است و مسئول جایگزینی سیاست‌ها، بروزرسانی‌ها یا آزمون امنیتی سرویس WhiteDNS نیست.

## محدوده پشتیبانی

- نصب‌کننده نسخه جاری DNSFabric و فایل‌های این مخزن تحت نگهداری پروژه‌اند.
- بسته WhiteDNS، موتورهای Soar/CottenDNS/MasterDNS و اپ اندروید دارای توسعه‌دهنده و مسیر گزارش جداگانه‌اند.
- این پروژه به‌طور پیش‌فرض API Token کلادفلر را دریافت نمی‌کند.
- روی سرورهای دارای پنل قبلی، نصب به‌صورت پیشگیرانه متوقف می‌شود.

## گزارش آسیب‌پذیری

**اطلاعات حساس را در GitHub Issues عمومی منتشر نکنید.** برای آسیب‌پذیری‌های امنیتی از <a href="https://github.com/Net-spoof/dnsfabric/security/advisories/new">Private Vulnerability Reporting / GitHub Security Advisory</a> در صورت فعال بودن این قابلیت استفاده کنید. در غیر این صورت پیش از ارسال جزئیات حساس از مسیر خصوصی مورد اطمینان با مالک مخزن هماهنگ کنید.

گزارش مفید شامل مراحل بازتولید، نسخه، اثر امنیتی و راهکار کاهش ریسک است. رمزها، توکن‌ها، کلیدهای خصوصی، دیتابیس و لینک خصوصی پنل را ضمیمه نکنید.

## نکات امنیتی نصب و نگهداری

- فقط مخزن رسمی و بسته بالادستی معتبر را استفاده کنید.
- قبل از آپدیت از اطلاعات پنل بکاپ بگیرید و بکاپ را خارج از VPS و به‌صورت محرمانه نگهداری کنید.
- دسترسی SSH، فایل‌های حساس و مسیر خصوصی ورود پنل را محدود کنید.
- checksum آرشیو، سلامت انتشار را بررسی می‌کند ولی جای ممیزی امنیتی و تأیید امضای مستقل انتشار نیست.
- HTTPS پنل و گواهی آن را دوره‌ای بررسی کنید.
- هیچ پیام «نصب موفق» نباید جایگزین آزمایش میدانی شبکه کاربران شود.

</div>

---

## Security disclosure (English)

Please **do not publish secrets or exploitable vulnerability details in public issues**. Use GitHub private vulnerability reporting if enabled for this repository. Otherwise contact the repository owner privately before sharing sensitive material. WhiteDNS and embedded engines are independent upstream projects and should receive upstream-specific security reports through their respective channels.
