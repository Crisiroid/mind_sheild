# مستندات بک‌اند اپلیکیشن روانشناسی

## معرفی پروژه

این پروژه یک API بک‌اند برای اپلیکیشن روانشناسی و سلامت روان است که با زبان Go نوشته شده. سیستم امکانات کاملی برای مدیریت کاربران، تمرین‌های روانشناسی، گزارش‌ها و داشبورد ادمین فراهم می‌کند.

## تکنولوژی‌ها

| تکنولوژی | نسخه | کاربرد |
|----------|------|--------|
| Go | 1.25.0 | زبان برنامه‌نویسی |
| Echo | 4.15.4 | فریم‌ورک HTTP و مدیریت روت‌ها |
| GORM | 1.31.1 | ORM برای عملیات دیتابیس |
| PostgreSQL Driver | 1.6.0 | اتصال به دیتابیس |
| golang-jwt | 5.3.1 | احراز هویت JWT |
| go-playground/validator | 10.30.3 | اعتبارسنجی داده‌ها |
| google/uuid | 1.6.0 | تولید و مدیریت UUID |
| godotenv | 1.5.1 | مدیریت متغیرهای محیطی |

## معماری پروژه

پروژه از الگوی Clean Architecture پیروی می‌کند و لایه‌های مختلفی دارد:

- **لایه Handler**: مدیریت درخواست‌ها و پاسخ‌های HTTP
- **لایه Service**: شامل منطق کسب‌وکار
- **لایه Repository**: ارتباط با دیتابیس
- **لایه Models**: تعریف ساختار داده‌ها و جداول دیتابیس
- **لایه Middleware**: وظایف مشترک مثل احراز هویت و لاگ کردن
- **پکیج Schemas**: شامل ثابت‌ها، مسیرهای روت و ساختار داده‌ها

## متغیرهای محیطی

| متغیر | توضیح | مقدار پیش‌فرض/مثال |
|-------|-------|-------------------|
| PORT | پورت سرور | 8080 |
| ENV | محیط اجرا | development یا production |
| DB_HOST | آدرس دیتابیس | localhost |
| DB_PORT | پورت دیتابیس | 5432 |
| DB_NAME | نام دیتابیس | psychology_db |
| DB_USER | نام کاربری دیتابیس | postgres |
| DB_PASSWORD | رمز عبور دیتابیس | (تنظیم شده) |
| JWT_SECRET | کلید امضای JWT | (تنظیم شده) |
| JWT_EXPIRY | مدت اعتبار توکن | 24h |
| UPLOAD_PATH | مسیر ذخیره فایل‌ها | ./uploads |
| BASE_URL | آدرس پایه سرویس | http://localhost:8080 |

## ساختار روت‌ها

روت‌ها به سه دسته اصلی تقسیم می‌شوند:

- **روت‌های عمومی**: بدون نیاز به احراز هویت
- **روت‌های کاربر**: نیاز به JWT معتبر با نقش user
- **روت‌های ادمین**: نیاز به JWT معتبر با نقش admin

## روت‌های عمومی

این روت‌ها نیاز به احراز هویت ندارند.

| متد | مسیر | توضیح |
|------|------|-------|
| GET | /health | بررسی سلامت سرویس با زمان فعلی |
| GET | /api/v1/public/health | بررسی سلامت عمومی سرویس |
| POST | /api/v1/auth/user/register | ثبت‌نام کاربر جدید |
| POST | /api/v1/auth/user/login | ورود کاربر |
| POST | /api/v1/auth/user/refresh | بروزرسانی توکن کاربر |
| POST | /api/v1/auth/admin/login | ورود ادمین |
| POST | /api/v1/auth/admin/refresh | بروزرسانی توکن ادمین |

## روت‌های احراز هویت کاربر

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/auth/user/logout | خروج کاربر |
| POST | /api/v1/auth/user/change-password | تغییر رمز عبور کاربر |

## روت‌های احراز هویت ادمین

نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/auth/admin/logout | خروج ادمین |
| POST | /api/v1/auth/admin/change-password | تغییر رمز عبور ادمین |

## روت‌های پروفایل کاربر

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| GET | /api/v1/users/me | دریافت پروفایل کاربر فعلی |
| PUT | /api/v1/users/me | بروزرسانی پروفایل کاربر فعلی |
| POST | /api/v1/users/me/sync | همگام‌سازی داده‌های کاربر |
| POST | /api/v1/users/me/agreement | پذیرش توافقنامه |

## روت‌های مدیریت کاربران (ادمین)

نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| GET | /api/v1/users | لیست تمام کاربران |
| POST | /api/v1/users | ایجاد کاربر جدید |
| GET | /api/v1/users/:id | دریافت اطلاعات کاربر بر اساس شناسه |
| PUT | /api/v1/users/:id | بروزرسانی اطلاعات کاربر |
| DELETE | /api/v1/users/:id | حذف کاربر |
| GET | /api/v1/users/by-phone | جستجوی کاربر بر اساس شماره تلفن |
| GET | /api/v1/users/stats | آمار کاربران |
| GET | /api/v1/users/activity-trend | روند فعالیت کاربران |
| GET | /api/v1/users/login-analytics | تحلیل ورود کاربران |
| GET | /api/v1/users/agreement-stats | آمار توافقنامه‌ها |
| GET | /api/v1/users/app-version-distribution | توزیع نسخه‌های اپلیکیشن |
| GET | /api/v1/users/inactive | کاربران غیرفعال |
| GET | /api/v1/users/engagement | میزان مشارکت کاربران |
| GET | /api/v1/users/export | خروجی اکسل از کاربران |
| POST | /api/v1/users/:id/accept-agreement | پذیرش توافقنامه برای کاربر |
| POST | /api/v1/users/update-login-info | بروزرسانی اطلاعات ورود |

## روت‌های تقویم روزانه

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/calendars | ایجاد ورودی تقویم جدید |
| GET | /api/v1/calendars | لیست ورودی‌های تقویم کاربر |
| GET | /api/v1/calendars/:id | دریافت ورودی تقویم بر اساس شناسه |
| PUT | /api/v1/calendars/:id | بروزرسانی ورودی تقویم |
| DELETE | /api/v1/calendars/:id | حذف ورودی تقویم |
| GET | /api/v1/calendars/stats/completion | آمار تکمیل تقویم |
| GET | /api/v1/calendars/stats/progress | پیشرفت در محدوده روزها |
| GET | /api/v1/calendars/stats/streak | تحلیل تداوم فعالیت |

## روت‌های مثلث احساسات

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/emotion-interactions | ایجاد تعامل احساسی جدید |
| GET | /api/v1/emotion-interactions | لیست تعاملات احساسی |
| GET | /api/v1/emotion-interactions/:id | دریافت تعامل احساسی بر اساس شناسه |
| PUT | /api/v1/emotion-interactions/:id | بروزرسانی تعامل احساسی |
| DELETE | /api/v1/emotion-interactions/:id | حذف تعامل احساسی |

## روت‌های رویدادهای استرس

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/stress-events | ایجاد رویداد استرس جدید |
| GET | /api/v1/stress-events | لیست رویدادهای استرس |
| GET | /api/v1/stress-events/:id | دریافت رویداد استرس بر اساس شناسه |
| PUT | /api/v1/stress-events/:id | بروزرسانی رویداد استرس |
| DELETE | /api/v1/stress-events/:id | حذف رویداد استرس |

## روت‌های نقشه تنش بدنی

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/body-tension-maps | ایجاد نقشه تنش بدنی جدید |
| GET | /api/v1/body-tension-maps | لیست نقشه‌های تنش بدنی |
| GET | /api/v1/body-tension-maps/:id | دریافت نقشه تنش بدنی بر اساس شناسه |
| PUT | /api/v1/body-tension-maps/:id | بروزرسانی نقشه تنش بدنی |
| DELETE | /api/v1/body-tension-maps/:id | حذف نقشه تنش بدنی |

## روت‌های تمرین تنفس

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/breathing-sessions | ایجاد جلسه تنفس جدید |
| GET | /api/v1/breathing-sessions | لیست جلسات تنفس |
| GET | /api/v1/breathing-sessions/:id | دریافت جلسه تنفس بر اساس شناسه |
| PUT | /api/v1/breathing-sessions/:id | بروزرسانی جلسه تنفس |
| DELETE | /api/v1/breathing-sessions/:id | حذف جلسه تنفس |

## روت‌های بازی شناختی

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/cognitive-games | ایجاد بازی شناختی جدید |
| GET | /api/v1/cognitive-games | لیست بازی‌های شناختی |
| GET | /api/v1/cognitive-games/:id | دریافت بازی شناختی بر اساس شناسه |
| PUT | /api/v1/cognitive-games/:id | بروزرسانی بازی شناختی |
| DELETE | /api/v1/cognitive-games/:id | حذف بازی شناختی |

## روت‌های بایدهای ذهنی

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/mental-musts | ایجاد باید ذهنی جدید |
| GET | /api/v1/mental-musts | لیست بایدهای ذهنی |
| GET | /api/v1/mental-musts/:id | دریافت باید ذهنی بر اساس شناسه |
| PUT | /api/v1/mental-musts/:id | بروزرسانی باید ذهنی |
| DELETE | /api/v1/mental-musts/:id | حذف باید ذهنی |

## روت‌های افکار منفی

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/negative-thoughts | ایجاد فکر منفی جدید |
| GET | /api/v1/negative-thoughts | لیست افکار منفی |
| GET | /api/v1/negative-thoughts/:id | دریافت فکر منفی بر اساس شناسه |
| PUT | /api/v1/negative-thoughts/:id | بروزرسانی فکر منفی |
| DELETE | /api/v1/negative-thoughts/:id | حذف فکر منفی |

## روت‌های دادگاه ذهن

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/mind-court-evidence | ایجاد شواهد دادگاه ذهن جدید |
| GET | /api/v1/mind-court-evidence | لیست شواهد دادگاه ذهن |
| GET | /api/v1/mind-court-evidence/:id | دریافت شواهد دادگاه ذهن بر اساس شناسه |
| PUT | /api/v1/mind-court-evidence/:id | بروزرسانی شواهد دادگاه ذهن |
| DELETE | /api/v1/mind-court-evidence/:id | حذف شواهد دادگاه ذهن |

## روت‌های تمرین تعارض

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/conflict-exercises | ایجاد تمرین تعارض جدید |
| GET | /api/v1/conflict-exercises | لیست تمرین‌های تعارض |
| GET | /api/v1/conflict-exercises/:id | دریافت تمرین تعارض بر اساس شناسه |
| PUT | /api/v1/conflict-exercises/:id | بروزرسانی تمرین تعارض |
| DELETE | /api/v1/conflict-exercises/:id | حذف تمرین تعارض |

## روت‌های ردیاب خلق‌و‌خو

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/mood-tracker | ایجاد ردیاب خلق‌و‌خو جدید |
| GET | /api/v1/mood-tracker | لیست ردیاب‌های خلق‌و‌خو |
| GET | /api/v1/mood-tracker/:id | دریافت ردیاب خلق‌و‌خو بر اساس شناسه |
| PUT | /api/v1/mood-tracker/:id | بروزرسانی ردیاب خلق‌و‌خو |
| DELETE | /api/v1/mood-tracker/:id | حذف ردیاب خلق‌و‌خو |

## روت‌های نقش‌ها و ارزش‌ها

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/roles-values | ایجاد نقش و ارزش جدید |
| GET | /api/v1/roles-values | لیست نقش‌ها و ارزش‌ها |
| GET | /api/v1/roles-values/:id | دریافت نقش و ارزش بر اساس شناسه |
| PUT | /api/v1/roles-values/:id | بروزرسانی نقش و ارزش |
| DELETE | /api/v1/roles-values/:id | حذف نقش و ارزش |

## روت‌های افکار آسمان

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/sky-thoughts | ایجاد فکر آسمان جدید |
| GET | /api/v1/sky-thoughts | لیست افکار آسمان |
| GET | /api/v1/sky-thoughts/:id | دریافت فکر آسمان بر اساس شناسه |
| PUT | /api/v1/sky-thoughts/:id | بروزرسانی فکر آسمان |
| DELETE | /api/v1/sky-thoughts/:id | حذف فکر آسمان |

## روت‌های ذهن‌آگاهی

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/mindful-timers | ایجاد تایمر ذهن‌آگاه جدید |
| GET | /api/v1/mindful-timers | لیست تایمرهای ذهن‌آگاه |
| GET | /api/v1/mindful-timers/:id | دریافت تایمر ذهن‌آگاه بر اساس شناسه |
| PUT | /api/v1/mindful-timers/:id | بروزرسانی تایمر ذهن‌آگاه |
| DELETE | /api/v1/mindful-timers/:id | حذف تایمر ذهن‌آگاه |
| POST | /api/v1/acceptance-exercises | ایجاد تمرین پذیرش جدید |
| GET | /api/v1/acceptance-exercises | لیست تمرین‌های پذیرش |
| GET | /api/v1/acceptance-exercises/:id | دریافت تمرین پذیرش بر اساس شناسه |
| PUT | /api/v1/acceptance-exercises/:id | بروزرسانی تمرین پذیرش |
| DELETE | /api/v1/acceptance-exercises/:id | حذف تمرین پذیرش |

## روت‌های گزارشات کاربر

نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/reports/weekly | ایجاد گزارش هفتگی جدید |
| GET | /api/v1/reports/weekly | لیست گزارش‌های هفتگی |
| GET | /api/v1/reports/weekly/:id | دریافت گزارش هفتگی بر اساس شناسه |
| PUT | /api/v1/reports/weekly/:id | بروزرسانی گزارش هفتگی |
| DELETE | /api/v1/reports/weekly/:id | حذف گزارش هفتگی |
| GET | /api/v1/reports/dashboard | داشبورد گزارشات |
| GET | /api/v1/reports/user-activity | فعالیت کاربر |
| GET | /api/v1/reports/stress-analytics | تحلیل استرس |
| GET | /api/v1/reports/body-tension | گزارش تنش بدنی |
| GET | /api/v1/reports/cognitive-patterns | الگوهای شناختی |
| GET | /api/v1/reports/mood-trends | روندهای خلق‌و‌خو |
| GET | /api/v1/reports/engagement | میزان مشارکت |
| GET | /api/v1/reports/weekly-progress | پیشرفت هفتگی |
| GET | /api/v1/reports/export | خروجی داده‌ها |
| GET | /api/v1/reports/weekly-stats | آمار هفتگی |

## روت‌های محتوای رسانه‌ای هفتگی

### روت‌های کاربر
نیاز به JWT معتبر با نقش user.

| متد | مسیر | توضیح |
|------|------|-------|
| GET | /api/v1/media/weekly | لیست محتوای رسانه‌ای |
| GET | /api/v1/media/weekly/:id | دریافت محتوا بر اساس شناسه |
| GET | /api/v1/media/weekly/by-week/:week_number | دریافت محتوا بر اساس شماره هفته |
| GET | /api/v1/media/weekly/:id/download | دانلود محتوا |

### روت‌های ادمین
نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/media/weekly | آپلود محتوای رسانه‌ای جدید |
| GET | /api/v1/media/weekly | لیست محتوای رسانه‌ای |
| GET | /api/v1/media/weekly/:id | دریافت محتوا بر اساس شناسه |
| PUT | /api/v1/media/weekly/:id | بروزرسانی محتوای رسانه‌ای |
| DELETE | /api/v1/media/weekly/:id | حذف محتوای رسانه‌ای |
| GET | /api/v1/media/weekly/by-week/:week_number | دریافت محتوا بر اساس شماره هفته |
| GET | /api/v1/media/weekly/:id/download | دانلود محتوا |

**مشخصات آپلود:**
- حداکثر حجم فایل: 100 مگابایت
- فرمت‌های پشتیبانی شده: mp3, wav, ogg, m4a, mp4, avi, mkv, pdf, doc, docx, jpg, png و سایر فرمت‌ها
- مسیر ذخیره‌سازی: ./uploads/weekly-media/

## روت‌های پروفایل ادمین

نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| GET | /api/v1/admin/me | دریافت پروفایل ادمین فعلی |
| PUT | /api/v1/admin/me | بروزرسانی پروفایل ادمین فعلی |

## روت‌های مدیریت ادمین‌ها

نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/admin/users | ایجاد ادمین جدید |
| GET | /api/v1/admin/users | لیست ادمین‌ها |
| GET | /api/v1/admin/users/:id | دریافت ادمین بر اساس شناسه |
| PUT | /api/v1/admin/users/:id | بروزرسانی اطلاعات ادمین |
| DELETE | /api/v1/admin/users/:id | حذف ادمین |
| POST | /api/v1/admin/users/:id/deactivate | غیرفعال کردن ادمین |

## روت‌های مدیریت نقش‌های ادمین

نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/admin/roles | ایجاد نقش جدید |
| GET | /api/v1/admin/roles | لیست نقش‌ها |
| GET | /api/v1/admin/roles/:id | دریافت نقش بر اساس شناسه |
| PUT | /api/v1/admin/roles/:id | بروزرسانی نقش |
| DELETE | /api/v1/admin/roles/:id | حذف نقش |

## روت‌های گزارشات ادمین

نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| POST | /api/v1/admin/reports | ایجاد گزارش کاربر جدید |
| GET | /api/v1/admin/reports | لیست گزارشات کاربران |
| GET | /api/v1/admin/reports/:id | دریافت گزارش بر اساس شناسه |
| DELETE | /api/v1/admin/reports/:id | حذف گزارش |
| GET | /api/v1/reports/dashboard | داشبورد گزارشات |
| GET | /api/v1/reports/user-activity | فعالیت کاربران |
| GET | /api/v1/reports/stress-analytics | تحلیل استرس |
| GET | /api/v1/reports/body-tension | گزارش تنش بدنی |
| GET | /api/v1/reports/cognitive-patterns | الگوهای شناختی |
| GET | /api/v1/reports/mood-trends | روندهای خلق‌و‌خو |
| GET | /api/v1/reports/engagement | میزان مشارکت |
| GET | /api/v1/reports/weekly-progress | پیشرفت هفتگی |
| GET | /api/v1/reports/export | خروجی داده‌ها |

## روت‌های لاگ سیستم

نیاز به JWT معتبر با نقش admin.

| متد | مسیر | توضیح |
|------|------|-------|
| GET | /api/v1/admin/logs | لیست لاگ‌های سیستم |
| GET | /api/v1/admin/logs/:id | دریافت لاگ بر اساس شناسه |

## روت‌های دسترسی به داده‌ها (ادمین)

نیاز به JWT معتبر با نقش admin. این روت‌ها برای مشاهده داده‌های تمام کاربران استفاده می‌شوند.

| متد | مسیر | توضیح |
|------|------|-------|
| GET | /api/v1/calendars | لیست تمام ورودی‌های تقویم |
| GET | /api/v1/emotion-interactions | لیست تمام تعاملات احساسی |
| GET | /api/v1/stress-events | لیست تمام رویدادهای استرس |
| GET | /api/v1/body-tension-maps | لیست تمام نقشه‌های تنش بدنی |
| GET | /api/v1/breathing-sessions | لیست تمام جلسات تنفس |
| GET | /api/v1/cognitive-games | لیست تمام بازی‌های شناختی |
| GET | /api/v1/mental-musts | لیست تمام بایدهای ذهنی |
| GET | /api/v1/negative-thoughts | لیست تمام افکار منفی |
| GET | /api/v1/mind-court-evidence | لیست تمام شواهد دادگاه ذهن |
| GET | /api/v1/conflict-exercises | لیست تمام تمرین‌های تعارض |
| GET | /api/v1/mood-tracker | لیست تمام ردیاب‌های خلق‌و‌خو |
| GET | /api/v1/roles-values | لیست تمام نقش‌ها و ارزش‌ها |
| GET | /api/v1/sky-thoughts | لیست تمام افکار آسمان |
| GET | /api/v1/mindful-timers | لیست تمام تایمرهای ذهن‌آگاه |
| GET | /api/v1/acceptance-exercises | لیست تمام تمرین‌های پذیرش |
| GET | /api/v1/reports/weekly | لیست تمام گزارشات هفتگی |

## مدل‌های داده

| مدل | فیلدهای اصلی |
|------|-------------|
| User | شناسه، شماره تلفن، نام، نام خانوادگی، تاریخ تولد، جنسیت، نسخه اپلیکیشن، آخرین ورود، وضعیت فعال بودن، وضعیت توافقنامه |
| Admin | شناسه، نام کاربری، رمز عبور hash شده، نام کامل، ایمیل، شناسه نقش، وضعیت فعال بودن |
| Admin Role | شناسه، نام نقش، توضیحات، مجوزها |
| Daily Calendar | شناسه، شناسه کاربر، تاریخ، تمرین‌های تکمیل شده، وضعیت تکمیل |
| Emotion Triangle | شناسه، شناسه کاربر، تاریخ، احساس اصلی، شدت، یادداشت |
| Stress Event | شناسه، شناسه کاربر، تاریخ، رویداد، شدت استرس، محل بدن، یادداشت |
| Body Tension Map | شناسه، شناسه کاربر، تاریخ، نواحی بدن، شدت تنش |
| Breathing Session | شناسه، شناسه کاربر، تاریخ، مدت زمان، نوع تمرین، وضعیت تکمیل |
| Cognitive Game | شناسه، شناسه کاربر، تاریخ، نوع بازی، امتیاز، سطح |
| Mental Must | شناسه، شناسه کاربر، تاریخ، متن باید، شدت، یادداشت |
| Negative Thought | شناسه، شناسه کاربر، تاریخ، متن فکر، نوع تحریف شناختی، شدت |
| Mind Court | شناسه، شناسه کاربر، تاریخ، فکر، شواهد موافق، شواهد مخالف، نتیجه |
| Conflict Exercise | شناسه، شناسه کاربر، تاریخ، موقعیت، احساسات، راه حل |
| Mood Tracker | شناسه، شناسه کاربر، تاریخ، خلق‌و‌خو، شدت، عوامل |
| Roles & Values | شناسه، شناسه کاربر، تاریخ، نقش، ارزش، میزان تطابق |
| Sky Thought | شناسه، شناسه کاربر، تاریخ، فکر، وضعیت رها کردن |
| Mindful Timer | شناسه، شناسه کاربر، تاریخ، مدت زمان، نوع مدیتیشن |
| Acceptance Exercise | شناسه، شناسه کاربر، تاریخ، موقعیت، سطح پذیرش، یادداشت |
| Weekly Report | شناسه، شناسه کاربر، شماره هفته، سال، آمار تمرین‌ها، خلاصه |
| Weekly Media | شناسه، شماره هفته، عنوان، توضیحات، مسیر فایل، نوع فایل، حجم فایل |
| System Log | شناسه، سطح لاگ، پیام، شناسه کاربر، زمان وقوع |

## نکات مهم

- تمام روت‌هایی که نیاز به احراز هویت دارند باید هدر Authorization با مقدار Bearer token ارسال کنند
- سیستم از دو نقش user و admin پشتیبانی می‌کند
- هر کاربر فقط به داده‌های خود دسترسی دارد مگر اینکه نقش admin داشته باشد
- ادمین‌ها به تمام داده‌های کاربران دسترسی دارند
- رمز عبور کاربران و ادمین‌ها با الگوریتم bcrypt hash می‌شود
- توکن JWT شامل شناسه کاربر و نقش است
- زمان انقضای توکن قابل تنظیم در فایل .env است
- فایل‌های آپلود شده در مسیر ./uploads ذخیره می‌شوند
- دیتابیس از PostgreSQL استفاده می‌کند
- ارتباط با دیتابیس از طریق ORM ابزار GORM انجام می‌شود
- تمام مدل‌ها شامل فیلدهای CreatedAt و UpdatedAt هستند
- شناسه‌ها از نوع UUID هستند

## نحوه اجرا

1. فایل .env را با مقادیر مناسب پیکربندی کنید
2. دستور `go mod tidy` را برای نصب وابستگی‌ها اجرا کنید
3. دستور `go run cmd/server/main.go` را برای اجرای سرور اجرا کنید
4. سرور روی پورت 8080 اجرا می‌شود
5. می‌توانید از ابزارهایی مثل Postman یا curl برای تست API استفاده کنید

## ساختار پوشه‌ها

| پوشه | توضیح |
|------|-------|
| cmd/server | فایل‌های اصلی شروع برنامه |
| internal/handler | handler های HTTP |
| internal/service | منطق کسب‌وکار |
| internal/repository | لایه دسترسی به داده |
| internal/models | مدل‌های داده |
| internal/middleware | middleware ها |
| internal/config | پیکربندی |
| internal/database | اتصال به دیتابیس |
| pkg/schemas | ساختار داده‌ها و ثابت‌ها |
| pkg/response | فرمت پاسخ‌ها |
| pkg/validator | توابع اعتبارسنجی |
| pkg/util | توابع کمکی |
| uploads | فایل‌های آپلود شده |

## پشتیبانی

برای گزارش مشکلات یا درخواست امکانات جدید با تیم توسعه تماس بگیرید.
