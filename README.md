# مستندات کامل بک اند اپلیکیشن روانشناسی

## معرفی پروژه

این پروژه یک API بک اند برای اپلیکیشن روانشناسی و سلامت روان است که با زبان برنامه نویسی Go نوشته شده. این سیستم امکانات کاملی برای مدیریت کاربران، تمرین های روانشناسی، گزارشات و داشبورد ادمین فراهم می کند.

## تکنولوژی های استفاده شده

زبان برنامه نویسی Go نسخه 1.25.0

فریم ورک Echo نسخه 4.15.4 برای مدیریت روت ها و HTTP handler ها

ابزار GORM نسخه 1.31.1 برای ارتباط با دیتابیس

دیتابیس PostgreSQL نسخه 1.6.0 درایور

احراز هویت JWT با کتابخانه golang-jwt نسخه 5.3.1

اعتبارسنجی داده ها با go-playground/validator نسخه 10.30.3

مدیریت UUID با google/uuid نسخه 1.6.0

مدیریت متغیرهای محیطی با godotenv نسخه 1.5.1

## معماری پروژه

پروژه از الگوی Clean Architecture پیروی می کند

لایه handler مسئول مدیریت درخواست های HTTP و پاسخ ها

لایه service شامل منطق کسب و کار

لایه repository مسئول ارتباط با دیتابیس

لایه models ساختار داده ها و جداول دیتابیس را تعریف می کند

لایه middleware وظایف مشترک مثل احراز هویت و لاگ کردن

پکیج schemas شامل ثابت ها، مسیرهای روت و ساختار داده ها

## متغیرهای محیطی

پورت پیش فرض 8080

محیط اجرا development یا production

تنظیمات دیتابیس شامل host، port، name، user، password

کلید JWT برای امضای توکن ها

مدت اعتبار توکن JWT پیش فرض 24 ساعت

مسیر ذخیره فایل های آپلود شده ./uploads

آدرس پایه سرویس BASE_URL

## ساختار روت ها

روت ها به سه دسته اصلی تقسیم می شوند

روت های عمومی بدون نیاز به احراز هویت

روت های کاربر با نقش user

روت های ادمین با نقش admin

## روت های عمومی

GET /health بررسی سلامت سرویس با زمان فعلی

GET /api/v1/public/health بررسی سلامت عمومی سرویس

POST /api/v1/auth/user/register ثبت نام کاربر جدید

POST /api/v1/auth/user/login ورود کاربر

POST /api/v1/auth/user/refresh بروزرسانی توکن کاربر

POST /api/v1/auth/admin/login ورود ادمین

POST /api/v1/auth/admin/refresh بروزرسانی توکن ادمین

## روت های احراز هویت کاربر

نیاز به JWT معتبر با نقش user

POST /api/v1/auth/user/logout خروج کاربر

POST /api/v1/auth/user/change-password تغییر رمز عبور کاربر

## روت های احراز هویت ادمین

نیاز به JWT معتبر با نقش admin

POST /api/v1/auth/admin/logout خروج ادمین

POST /api/v1/auth/admin/change-password تغییر رمز عبور ادمین

## روت های پروفایل کاربر

نیاز به JWT معتبر با نقش user

GET /api/v1/users/me دریافت پروفایل کاربر فعلی

PUT /api/v1/users/me بروزرسانی پروفایل کاربر فعلی

POST /api/v1/users/me/sync همگام سازی داده های کاربر

POST /api/v1/users/me/agreement پذیرش توافقنامه

## روت های مدیریت کاربران توسط ادمین

نیاز به JWT معتبر با نقش admin

GET /api/v1/users لیست تمام کاربران

POST /api/v1/users ایجاد کاربر جدید

GET /api/v1/users/:id دریافت اطلاعات کاربر بر اساس شناسه

PUT /api/v1/users/:id بروزرسانی اطلاعات کاربر

DELETE /api/v1/users/:id حذف کاربر

GET /api/v1/users/by-phone جستجوی کاربر بر اساس شماره تلفن

GET /api/v1/users/stats آمار کاربران

GET /api/v1/users/activity-trend روند فعالیت کاربران

GET /api/v1/users/login-analytics تحلیل ورود کاربران

GET /api/v1/users/agreement-stats آمار توافقنامه ها

GET /api/v1/users/app-version-distribution توزیع نسخه های اپلیکیشن

GET /api/v1/users/inactive کاربران غیرفعال

GET /api/v1/users/engagement میزان مشارکت کاربران

GET /api/v1/users/export خروجی اکسل از کاربران

POST /api/v1/users/:id/accept-agreement پذیرش توافقنامه برای کاربر

POST /api/v1/users/update-login-info بروزرسانی اطلاعات ورود

## روت های تقویم روزانه

نیاز به JWT معتبر با نقش user

POST /api/v1/calendars ایجاد ورودی تقویم جدید

GET /api/v1/calendars لیست ورودی های تقویم کاربر

GET /api/v1/calendars/:id دریافت ورودی تقویم بر اساس شناسه

PUT /api/v1/calendars/:id بروزرسانی ورودی تقویم

DELETE /api/v1/calendars/:id حذف ورودی تقویم

GET /api/v1/calendars/stats/completion آمار تکمیل تقویم

GET /api/v1/calendars/stats/progress پیشرفت در محدوده روزها

GET /api/v1/calendars/stats/streak تحلیل تداوم فعالیت

## روت های مثلث احساسات

نیاز به JWT معتبر با نقش user

POST /api/v1/emotion-interactions ایجاد تعامل احساسی جدید

GET /api/v1/emotion-interactions لیست تعاملات احساسی

GET /api/v1/emotion-interactions/:id دریافت تعامل احساسی بر اساس شناسه

PUT /api/v1/emotion-interactions/:id بروزرسانی تعامل احساسی

DELETE /api/v1/emotion-interactions/:id حذف تعامل احساسی

## روت های رویدادهای استرس

نیاز به JWT معتبر با نقش user

POST /api/v1/stress-events ایجاد رویداد استرس جدید

GET /api/v1/stress-events لیست رویدادهای استرس

GET /api/v1/stress-events/:id دریافت رویداد استرس بر اساس شناسه

PUT /api/v1/stress-events/:id بروزرسانی رویداد استرس

DELETE /api/v1/stress-events/:id حذف رویداد استرس

## روت های نقشه تنش بدنی

نیاز به JWT معتبر با نقش user

POST /api/v1/body-tension-maps ایجاد نقشه تنش بدنی جدید

GET /api/v1/body-tension-maps لیست نقشه های تنش بدنی

GET /api/v1/body-tension-maps/:id دریافت نقشه تنش بدنی بر اساس شناسه

PUT /api/v1/body-tension-maps/:id بروزرسانی نقشه تنش بدنی

DELETE /api/v1/body-tension-maps/:id حذف نقشه تنش بدنی

## روت های تمرین تنفس

نیاز به JWT معتبر با نقش user

POST /api/v1/breathing-sessions ایجاد جلسه تنفس جدید

GET /api/v1/breathing-sessions لیست جلسات تنفس

GET /api/v1/breathing-sessions/:id دریافت جلسه تنفس بر اساس شناسه

PUT /api/v1/breathing-sessions/:id بروزرسانی جلسه تنفس

DELETE /api/v1/breathing-sessions/:id حذف جلسه تنفس

## روت های بازی شناختی

نیاز به JWT معتبر با نقش user

POST /api/v1/cognitive-games ایجاد بازی شناختی جدید

GET /api/v1/cognitive-games لیست بازی های شناختی

GET /api/v1/cognitive-games/:id دریافت بازی شناختی بر اساس شناسه

PUT /api/v1/cognitive-games/:id بروزرسانی بازی شناختی

DELETE /api/v1/cognitive-games/:id حذف بازی شناختی

## روت های باید های ذهنی

نیاز به JWT معتبر با نقش user

POST /api/v1/mental-musts ایجاد باید ذهنی جدید

GET /api/v1/mental-musts لیست باید های ذهنی

GET /api/v1/mental-musts/:id دریافت باید ذهنی بر اساس شناسه

PUT /api/v1/mental-musts/:id بروزرسانی باید ذهنی

DELETE /api/v1/mental-musts/:id حذف باید ذهنی

## روت های افکار منفی

نیاز به JWT معتبر با نقش user

POST /api/v1/negative-thoughts ایجاد فکر منفی جدید

GET /api/v1/negative-thoughts لیست افکار منفی

GET /api/v1/negative-thoughts/:id دریافت فکر منفی بر اساس شناسه

PUT /api/v1/negative-thoughts/:id بروزرسانی فکر منفی

DELETE /api/v1/negative-thoughts/:id حذف فکر منفی

## روت های دادگاه ذهن

نیاز به JWT معتبر با نقش user

POST /api/v1/mind-court-evidence ایجاد شواهد دادگاه ذهن جدید

GET /api/v1/mind-court-evidence لیست شواهد دادگاه ذهن

GET /api/v1/mind-court-evidence/:id دریافت شواهد دادگاه ذهن بر اساس شناسه

PUT /api/v1/mind-court-evidence/:id بروزرسانی شواهد دادگاه ذهن

DELETE /api/v1/mind-court-evidence/:id حذف شواهد دادگاه ذهن

## روت های تمرین تعارض

نیاز به JWT معتبر با نقش user

POST /api/v1/conflict-exercises ایجاد تمرین تعارض جدید

GET /api/v1/conflict-exercises لیست تمرین های تعارض

GET /api/v1/conflict-exercises/:id دریافت تمرین تعارض بر اساس شناسه

PUT /api/v1/conflict-exercises/:id بروزرسانی تمرین تعارض

DELETE /api/v1/conflict-exercises/:id حذف تمرین تعارض

## روت های ردیاب خلق و خو

نیاز به JWT معتبر با نقش user

POST /api/v1/mood-tracker ایجاد ردیاب خلق و خو جدید

GET /api/v1/mood-tracker لیست ردیاب های خلق و خو

GET /api/v1/mood-tracker/:id دریافت ردیاب خلق و خو بر اساس شناسه

PUT /api/v1/mood-tracker/:id بروزرسانی ردیاب خلق و خو

DELETE /api/v1/mood-tracker/:id حذف ردیاب خلق و خو

## روت های نقش ها و ارزش ها

نیاز به JWT معتبر با نقش user

POST /api/v1/roles-values ایجاد نقش و ارزش جدید

GET /api/v1/roles-values لیست نقش ها و ارزش ها

GET /api/v1/roles-values/:id دریافت نقش و ارزش بر اساس شناسه

PUT /api/v1/roles-values/:id بروزرسانی نقش و ارزش

DELETE /api/v1/roles-values/:id حذف نقش و ارزش

## روت های افکار آسمان

نیاز به JWT معتبر با نقش user

POST /api/v1/sky-thoughts ایجاد فکر آسمان جدید

GET /api/v1/sky-thoughts لیست افکار آسمان

GET /api/v1/sky-thoughts/:id دریافت فکر آسمان بر اساس شناسه

PUT /api/v1/sky-thoughts/:id بروزرسانی فکر آسمان

DELETE /api/v1/sky-thoughts/:id حذف فکر آسمان

## روت های ذهن آگاهی

نیاز به JWT معتبر با نقش user

POST /api/v1/mindful-timers ایجاد تایمر ذهن آگاه جدید

GET /api/v1/mindful-timers لیست تایمرهای ذهن آگاه

GET /api/v1/mindful-timers/:id دریافت تایمر ذهن آگاه بر اساس شناسه

PUT /api/v1/mindful-timers/:id بروزرسانی تایمر ذهن آگاه

DELETE /api/v1/mindful-timers/:id حذف تایمر ذهن آگاه

POST /api/v1/acceptance-exercises ایجاد تمرین پذیرش جدید

GET /api/v1/acceptance-exercises لیست تمرین های پذیرش

GET /api/v1/acceptance-exercises/:id دریافت تمرین پذیرش بر اساس شناسه

PUT /api/v1/acceptance-exercises/:id بروزرسانی تمرین پذیرش

DELETE /api/v1/acceptance-exercises/:id حذف تمرین پذیرش

## روت های گزارشات کاربر

نیاز به JWT معتبر با نقش user

POST /api/v1/reports/weekly ایجاد گزارش هفتگی جدید

GET /api/v1/reports/weekly لیست گزارش های هفتگی

GET /api/v1/reports/weekly/:id دریافت گزارش هفتگی بر اساس شناسه

PUT /api/v1/reports/weekly/:id بروزرسانی گزارش هفتگی

DELETE /api/v1/reports/weekly/:id حذف گزارش هفتگی

GET /api/v1/reports/dashboard داشبورد گزارشات

GET /api/v1/reports/user-activity فعالیت کاربر

GET /api/v1/reports/stress-analytics تحلیل استرس

GET /api/v1/reports/body-tension گزارش تنش بدنی

GET /api/v1/reports/cognitive-patterns الگوهای شناختی

GET /api/v1/reports/mood-trends روندهای خلق و خو

GET /api/v1/reports/engagement میزان مشارکت

GET /api/v1/reports/weekly-progress پیشرفت هفتگی

GET /api/v1/reports/export خروجی داده ها

GET /api/v1/reports/weekly-stats آمار هفتگی

## روت های محتوای رسانه ای هفتگی

روت های کاربر نیاز به JWT معتبر با نقش user

GET /api/v1/media/weekly لیست محتوای رسانه ای

GET /api/v1/media/weekly/:id دریافت محتوا بر اساس شناسه

GET /api/v1/media/weekly/by-week/:week_number دریافت محتوا بر اساس شماره هفته

GET /api/v1/media/weekly/:id/download دانلود محتوا

روت های ادمین نیاز به JWT معتبر با نقش admin

POST /api/v1/media/weekly آپلود محتوای رسانه ای جدید

GET /api/v1/media/weekly لیست محتوای رسانه ای

GET /api/v1/media/weekly/:id دریافت محتوا بر اساس شناسه

PUT /api/v1/media/weekly/:id بروزرسانی محتوای رسانه ای

DELETE /api/v1/media/weekly/:id حذف محتوای رسانه ای

GET /api/v1/media/weekly/by-week/:week_number دریافت محتوا بر اساس شماره هفته

GET /api/v1/media/weekly/:id/download دانلود محتوا

حداکثر حجم فایل 100 مگابایت

فرمت های پشتیبانی شده شامل mp3، wav، ogg، m4a، mp4، avi، mkv، pdf، doc، docx، jpg، png و سایر فرمت ها

فایل ها در مسیر ./uploads/weekly-media/ ذخیره می شوند

## روت های پروفایل ادمین

نیاز به JWT معتبر با نقش admin

GET /api/v1/admin/me دریافت پروفایل ادمین فعلی

PUT /api/v1/admin/me بروزرسانی پروفایل ادمین فعلی

## روت های مدیریت ادمین ها

نیاز به JWT معتبر با نقش admin

POST /api/v1/admin/users ایجاد ادمین جدید

GET /api/v1/admin/users لیست ادمین ها

GET /api/v1/admin/users/:id دریافت ادمین بر اساس شناسه

PUT /api/v1/admin/users/:id بروزرسانی اطلاعات ادمین

DELETE /api/v1/admin/users/:id حذف ادمین

POST /api/v1/admin/users/:id/deactivate غیرفعال کردن ادمین

## روت های مدیریت نقش های ادمین

نیاز به JWT معتبر با نقش admin

POST /api/v1/admin/roles ایجاد نقش جدید

GET /api/v1/admin/roles لیست نقش ها

GET /api/v1/admin/roles/:id دریافت نقش بر اساس شناسه

PUT /api/v1/admin/roles/:id بروزرسانی نقش

DELETE /api/v1/admin/roles/:id حذف نقش

## روت های گزارشات ادمین

نیاز به JWT معتبر با نقش admin

POST /api/v1/admin/reports ایجاد گزارش کاربر جدید

GET /api/v1/admin/reports لیست گزارشات کاربران

GET /api/v1/admin/reports/:id دریافت گزارش بر اساس شناسه

DELETE /api/v1/admin/reports/:id حذف گزارش

GET /api/v1/reports/dashboard داشبورد گزارشات

GET /api/v1/reports/user-activity فعالیت کاربران

GET /api/v1/reports/stress-analytics تحلیل استرس

GET /api/v1/reports/body-tension گزارش تنش بدنی

GET /api/v1/reports/cognitive-patterns الگوهای شناختی

GET /api/v1/reports/mood-trends روندهای خلق و خو

GET /api/v1/reports/engagement میزان مشارکت

GET /api/v1/reports/weekly-progress پیشرفت هفتگی

GET /api/v1/reports/export خروجی داده ها

## روت های لاگ سیستم

نیاز به JWT معتبر با نقش admin

GET /api/v1/admin/logs لیست لاگ های سیستم

GET /api/v1/admin/logs/:id دریافت لاگ بر اساس شناسه

## روت های داده های ادمین

نیاز به JWT معتبر با نقش admin

این روت ها برای مشاهده داده های تمام کاربران استفاده می شوند

GET /api/v1/calendars لیست تمام ورودی های تقویم

GET /api/v1/emotion-interactions لیست تمام تعاملات احساسی

GET /api/v1/stress-events لیست تمام رویدادهای استرس

GET /api/v1/body-tension-maps لیست تمام نقشه های تنش بدنی

GET /api/v1/breathing-sessions لیست تمام جلسات تنفس

GET /api/v1/cognitive-games لیست تمام بازی های شناختی

GET /api/v1/mental-musts لیست تمام باید های ذهنی

GET /api/v1/negative-thoughts لیست تمام افکار منفی

GET /api/v1/mind-court-evidence لیست تمام شواهد دادگاه ذهن

GET /api/v1/conflict-exercises لیست تمام تمرین های تعارض

GET /api/v1/mood-tracker لیست تمام ردیاب های خلق و خو

GET /api/v1/roles-values لیست تمام نقش ها و ارزش ها

GET /api/v1/sky-thoughts لیست تمام افکار آسمان

GET /api/v1/mindful-timers لیست تمام تایمرهای ذهن آگاه

GET /api/v1/acceptance-exercises لیست تمام تمرین های پذیرش

GET /api/v1/reports/weekly لیست تمام گزارشات هفتگی

## مدل های داده

کاربر شامل شناسه، شماره تلفن، نام، نام خانوادگی، تاریخ تولد، جنسیت، نسخه اپلیکیشن، آخرین زمان ورود، وضعیت فعال بودن، توافقنامه

ادمین شامل شناسه، نام کاربری، رمز عبور hash شده، نام کامل، ایمیل، شناسه نقش، وضعیت فعال بودن

نقش ادمین شامل شناسه، نام نقش، توضیحات، مجوزها

تقویم روزانه شامل شناسه، شناسه کاربر، تاریخ، لیست تمرین های تکمیل شده، وضعیت تکمیل

مثلث احساسات شامل شناسه، شناسه کاربر، تاریخ، احساس اصلی، شدت، یادداشت

رویداد استرس شامل شناسه، شناسه کاربر، تاریخ، رویداد، شدت استرس، محل بدن، یادداشت

نقشه تنش بدنی شامل شناسه، شناسه کاربر، تاریخ، نواحی بدن، شدت تنش

جلسه تنفس شامل شناسه، شناسه کاربر، تاریخ، مدت زمان، نوع تمرین، وضعیت تکمیل

بازی شناختی شامل شناسه، شناسه کاربر، تاریخ، نوع بازی، امتیاز، سطح

باید ذهنی شامل شناسه، شناسه کاربر، تاریخ، متن باید، شدت، یادداشت

فکر منفی شامل شناسه، شناسه کاربر، تاریخ، متن فکر، نوع تحریف شناختی، شدت

دادگاه ذهن شامل شناسه، شناسه کاربر، تاریخ، فکر، شواهد موافق، شواهد مخالف، نتیجه

تمرین تعارض شامل شناسه، شناسه کاربر، تاریخ، موقعیت، احساسات، راه حل

ردیاب خلق و خو شامل شناسه، شناسه کاربر، تاریخ، خلق و خو، شدت، عوامل

نقش و ارزش شامل شناسه، شناسه کاربر، تاریخ، نقش، ارزش، میزان تطابق

فکر آسمان شامل شناسه، شناسه کاربر، تاریخ، فکر، وضعیت رها کردن

تایمر ذهن آگاه شامل شناسه، شناسه کاربر، تاریخ، مدت زمان، نوع مدیتیشن

تمرین پذیرش شامل شناسه، شناسه کاربر، تاریخ، موقعیت، سطح پذیرش، یادداشت

گزارش هفتگی شامل شناسه، شناسه کاربر، شماره هفته، سال، آمار تمرین ها، خلاصه

محتوای رسانه ای هفتگی شامل شناسه، شماره هفته، عنوان، توضیحات، مسیر فایل، نوع فایل، حجم فایل

لاگ سیستم شامل شناسه، سطح لاگ، پیام، شناسه کاربر، زمان وقوع

## نکات مهم

تمام روت هایی که نیاز به احراز هویت دارند باید هدر Authorization با مقدار Bearer token ارسال کنند

سیستم از دو نقش user و admin پشتیبانی می کند

هر کاربر فقط به داده های خود دسترسی دارد مگر اینکه نقش admin داشته باشد

ادمین ها به تمام داده های کاربران دسترسی دارند

رمز عبور کاربران و ادمین ها با الگوریتم bcrypt hash می شود

توکن JWT شامل شناسه کاربر و نقش است

زمان انقضای توکن قابل تنظیم در فایل .env است

فایل های آپلود شده در مسیر ./uploads ذخیره می شوند

دیتابیس از PostgreSQL استفاده می کند

ارتباط با دیتابیس از طریق ORM ابزار GORM انجام می شود

تمام مدل ها شامل فیلدهای CreatedAt و UpdatedAt هستند

شناسه ها از نوع UUID هستند

## نحوه اجرا

ابتدا فایل .env را با مقادیر مناسب پیکربندی کنید

دستور go mod tidy را برای نصب وابستگی ها اجرا کنید

دستور go run cmd/server/main.go را برای اجرای سرور اجرا کنید

سرور روی پورت 8080 اجرا می شود

می توانید از ابزارهایی مثل Postman یا curl برای تست API استفاده کنید

## ساختار پوشه ها

cmd/server شامل فایل های اصلی شروع برنامه

internal/handler شامل handler های HTTP

internal/service شامل منطق کسب و کار

internal/repository شامل لایه دسترسی به داده

internal/models شامل مدل های داده

internal/middleware شامل middleware ها

internal/config شامل پیکربندی

internal/database شامل اتصال به دیتابیس

pkg/schemas شامل ساختار داده ها و ثابت ها

pkg/response شامل فرمت پاسخ ها

pkg/validator شامل توابع اعتبارسنجی

pkg/util شامل توابع کمکی

uploads شامل فایل های آپلود شده

## پشتیبانی

برای گزارش مشکلات یا درخواست امکانات جدید با تیم توسعه تماس بگیرید
