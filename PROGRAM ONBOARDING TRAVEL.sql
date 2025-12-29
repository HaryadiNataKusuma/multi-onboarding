PROGRAM ONBOARDING TRAVEL
INSURANCES : Input one by one oleh user
REGION : 
---

----SECTION REGIONS
di section Regions tolong ganti field Country menjadi button aja bro. ketika user klik tombol button itu maka akan muncul pop up yang isinya data dari table travel_service_development.countries. 
di dalam pop up tersebut user bisa pilih negara-negara mana saja yang akan masuk ke dalam Regions tersebut. selain itu tambahkan juga satu field di dalam pop up untuk menampung country id yang dipilih oleh user.
format data di dalam field country id adalah json, contoh : [94], jika user pilih country lebih dari satu maka formatnya sepert ini contohnya : [146,1,51,45,95,96,117,213,194,199,46,159,90]. dan user bisa menghapus country id ini jika terjadi salah pilih country.
format baku untuk Nama Regions
'INSURANCE CODE' + 'REGION TYPE'
contoh : DAMAI DOMESTIC, DAMAI ASEAN, DAMAI WORLDWIDE, DAMAI WORLDWIDE EXLUDE USCA, DAMAI ASIA PACIFIC.
--
dan ketika gue klik satu negara, seharusnya hanya negara itu saja yang terpilih, jangan semua negara. dan id negara tersebut langsung te-reflect pada field Country Ids.
----SECTION PRODCTS
di section products tolong ubah field : 
1. Insurance Name menjadi Type (dropdown : COUPLE, INDIVIDUAL, FAMILY)
3. Vehicle Type menjadi Region ID (dropdown : data ambil dari kolom name pada table travel_service_development.regions yang diinput oleh user)
4. hilangkan field Is Electric Vehicle
5. hilangkan field logo, tapi logo otomatis ambil dari kolom logo pada table table travel_service_development.insurances
TV-DAMAI-DOM-IND-01
DAMAI Domestic Deluxe Individu
TV-DAMAI-DOM-FAM-01
DAMAI Domestic Deluxe Family
sesuaikan kolom-kolom pada draft data section products dan RESULT section products dengan kolom-kolom yang ada di table travel_service_development.products. tapi tidak perlu menampilkan kolom id, protection_detail, how_to_claim, created_by, updated_by, updated_at, dan partner
pada draft data dan RESULT section products :
1. tambahkan field Insurance Type (Dropdown : DOMESTIC, WORLDWIDE, ASIA) dan data yang dipilih oleh user akan masuk ke kolom insurance_type pada table travel_service_delopment.products.
2. data dari field type yang dipilih oleh user akan masuk ke kolom type pada table travel_service_delopment.products.
3. kolom regiond id akan terisi oleh id dari table travel_service_delopment.regions berdasarkan Nama Region yang dipilih oleh user pada field Region
4. jika user memilih type INDIVIDUAL : kolom min_adult dan max_adult akan otomatis terisi dengan data angka 1, min_child dan max_child akan otomatis terisi dengan angka -1. jika user memilih type COUPLE : kolom min_adult akan otomatis terisi dengan data angka 1 dan max_adult akan otomatis terisi dengan data angka 2, min_child dan max_child akan otomatis terisi dengan angka -1. jika user memilih type FAMILY : kolom min_adult akan otomatis terisi dengan data angka 1 dan max_adult akan otomatis terisi dengan data angka 5, min_child akan otomatis terisi dengan angka 0 dan max_child akan otomatis terisi dengan angka 5.
5. kolom SUMMARY akan otomatis terisi dengan data SUMMARY_ (ditambah dengan code product yang simbol "-" sudah diconvert menjadi "_": SUMMARY_TV_DAMAI_DOM_IND_01)
6. kolom INSURANCE_DETAIL akan otomatis terisi dengan data INSURANCE_DETAIL_ (ditambah dengan code product yang simbol "-" sudah diconvert menjadi "_": INSURANCE_DETAIL_TV_DAMAI_DOM_IND_01)
7. kolom PROTECTION_DETAIL akan otomatis terisi dengan data PROTECTION_DETAIL_ (ditambah dengan code product yang simbol "-" sudah diconvert menjadi "_": PROTECTION_DETAIL_TV_DAMAI_DOM_IND_01)
8. kolom HOW_TO_CLAIM akan otomatis terisi dengan data HOW_TO_CLAIM_ (ditambah dengan code product yang simbol "-" sudah diconvert menjadi "_": HOW_TO_CLAIM_TV_DAMAI_DOM_IND_01)
GENERATED NAME : 
khusus untuk product yang protection type = ANNUAL maka program akan otomatis generate 2 product code. 
* 1 product code untuk perlindungan 90 Days (contoh : TV-DAMAI-AP-IND-ANN-01), 1 product code untuk perlindungan 180 Days (contoh : TV-DAMAI-AP-IND-ANN-01).
* nama pada kolom name table travel_service_delopment.products tambahkan kata (90 Days) pada product name yang perlindungan 90 Days
* nama pada kolom name table travel_service_delopment.products tambahkan kata (180 Days) pada product name yang perlindungan 180 Days
--product rules
untuk product_code yang premium_type nya = ANNUAL isi pada template excel file adalah sebagai berikut : 
product_code 90 days : start_days otomatis terisi dengan data angka 1 dan end_days otomatis terisi dengan data angka 90
product_code 180 days : start_days otomatis terisi dengan data angka 1 dan end_days otomatis terisi dengan data angka 180
--
eh bro, gue punya ide untuk gabungin section products dengan section addons. jadi field addon name (beserta fungsinya) yang ada di section addons lo tambahin di section products bro, tidak perlu field product codes dan field status (status otomatis terisi dengan data angka 1). jadi user tinggal klik tombol addons yang ada di section product untuk pilih addons apa yang menjadi perlindungan tambahan di product tersebut.
output dari addons yang sudah dipilih user langsung masuk ke database ke table travel_service_development.addons dan travel_service_development.addon_rules
--addons
jika gue pilih addons dari section products maka addons tersebut akan masuk ke table travel_service_development.addons, travel_service_development.addon_rules, dan travel_service_development.insurance_product_addon_mappings :
1. travel_service_development.addons : 
* kolom code isi dengan kode addons yang dipilih oleh user
* kolom name isi dengan nama addons yang diinput oleh user
* kolom name_my isi dengan nama addons yang diinput oleh user
* kolom name_en isi dengan nama addons yang diinput oleh user
* kolom is_active otomatis akan terisi dengan data angka 1
* kolom created_by otomatis akan terisi dengan data angka 1
* kolom created_at otomatis akan terisi dengan CURRENT_TIMESTAMP
* kolom updated_by otomati akan terisi dengan data angka 1 jika RESULT di edit
* kolom updated_by otomati akan terisi dengan CURRENT_TIMESTAMP jika RESULT di edit
2. travel_service_development.addon_rules : 
* kolom addon_code isi dengan kode addons yang dipilih oleh user
* kolom product_code isi dengan generated product code
* kolom insurance_code isi dengan insurance code
* kolom otomatis akan terisi dengan data {}
* kolom created_by otomatis akan terisi dengan data angka 1
* kolom created_at otomatis akan terisi dengan CURRENT_TIMESTAMP
* kolom updated_by otomati akan terisi dengan data angka 1 jika RESULT di edit
* kolom updated_by otomati akan terisi dengan CURRENT_TIMESTAMP jika RESULT di edit
3. travel_service_development.insurance_product_addon_mappings : 
* kolom insurance_code isi dengan insurance code
* kolom product_code isi dengan generated product code
* kolom addon_code isi dengan kode addons yang dipilih oleh user
* kolom created_by otomatis akan terisi dengan data angka 1
* kolom created_at otomatis akan terisi dengan CURRENT_TIMESTAMP
* kolom updated_by otomati akan terisi dengan data angka 1 jika RESULT di edit
* kolom updated_by otomati akan terisi dengan CURRENT_TIMESTAMP jika RESULT di edit
ketika user input product baru dan menambahkan addons, pasti akan terjadi duplicate pada table travel_service_development.addons. kalau addon sudah ada di table travel_service_development.addons maka lo gak perlu insert data. tapi kalau addon nya belum ada di travel_service_development.addons, maka perlu insert data.

addon_rule_id, start_condition, end_condition, value_type, value, duration_rule_type, min_adult, max_adult, max_age



0968a5cb-dc83-41c2-b230-d4d9dbd2d393
f0745170-a305-47f6-8e11-4b82a45f1925

TV-{ditambah dengan code insurances}-{ditambah dengan name Regions yang dipilih user (kalau name like '%DOMESTIC%' maka cukup kasih kode DOM, kalau name like '%WORLDWIDE%' maka cukup kasih kode WW, kalau name like '%EXCLUDE%' maka cukup kasih kode WWEX, kalau name like '%ASEAN%' maka cukup kasih kode ASN, kalau name like '%asia pacific%' maka cukup kasih kode AP)}-{ditambah dengan type yang dipilih user (kalau COUPLE maka cukup kasih kode COU, kalau INDIVIDUAL maka cukup kasih kode IND, kalau FAMILY maka cukup kasih kode FAM)}-{ditambah dengan kode 01(increment)}

STARTER,
SILVER,
GOLD,PLATINUM,CORPORATE,CORPORATE2,DIAMOND,PARTNER_MOLADIN,CORPORATE4,GREEN,TIED,CORPORATE3,QOALA,SKB,SHOPDRIVE,DIRECTUSEDCAR,DIRECTPROPERTY
commissions_service_local.commissions
agent_service_local.default_config_products
quotation_service_local.plan_commissions

untuk duration_tipe DAILAY
start_days : 1,5,7,9,11,16,21,26
end_days : 4,6,8,10,15,20,25,31
untuk duration_type PER_EXTRA_SEVEN_DAYS
start_days : 32
end_days : 180
untuk duration_type ANNUAL
start_days : 1
end_days : 365



'MV-CAR-ACA-03','MV-CAR-ACA-04','MV-CAR-ACA-05'

GREEN
SILVER
GOLD
DIAMOND
PLATINUM
CORPORATE
CORPORATE2
CORPORATE3
CORPORATE4
PARTNER_MOLADIN
TIED
QOALA
SKB
DIRECTPROPERTY

ADDON CODE - ADDON NAME :
COVID-PROTECTION - Covid Protection
TRAVEL-IBADAH - Travel Ibadah
WINTER-SPORT - Winter Sports
TRAVEL-PENDIDIKAN - Travel Pendidikan
TRIP-CANCELLATION-EXTENSION - Trip Cancellation Extension
GOLF-BENEFIT - Golf Benefit
LOSS-OF-HOTEL-RESERVATION - Loss of Hotel Reservation
PET-INSURANCE - Pet insurance
MISS-EVENT - Missed event
ADVENTURE - Adventures Activities Cover 
VISA-PROTECTION - Visa Protection
AMATEUR-SPORT-EVENT - Amateur Sports Event
BUSINESS-TRAVEL - Business Travel
SNOW-SPORT - Snow Sport Cover