package constant

const (
	HELP_MENU_FORMAT = `
*── 「 WELCOME 」 ──*

======================
➸ *Nama*: {{ pushname }}
➸ *Waktu*: {{ date }}
======================

Menu yang tersedia:

{{ menu }}
======================
➸ *Prefix*:
[!][#][?][.][@][&][*][-][=][+]
======================

Ketik *{{ prefix }}menu* _angka_index_ *atau* _nama_menu_
untuk membuka menu page yang dipilih.

Catatan:
Bot ini terdapat cooldown command selama *5 detik* setiap kali pemakaian.
`

	Y2MATE_DESC = `
*── 「 YOUTUBE 」 ──*

======================
*Url*: [url]
*Title*: [title]
*Duration*: [duration]
*Channel*: [channel]
======================

Available Downloads:

[downloads]
======================
Balas *MediaID* sebuah format untuk mendownload, contoh:
*auto* untuk memilih format *auto*
======================
`
	Y2MATE_RESULT_CAPTION = `
*── 「 YOUTUBE 」 ──*

======================
*Url*: [url]
*Title*: [title]
*Duration*: [duration]
*Channel*: [channel]
*Size*: [size]
*Format*: [format]
======================
`
	SNAPTIK_LIST = `
*── 「 TIKTOK 」 ──*

======================
*Username*: [username]
*Description*: [description]
======================

Available Slide:

[slides]
======================
Balas nomor sebuah slide untuk mendownload, contoh:
*1* untuk memilih slide pertama
======================
`
	SNAPTIK_RESULT = `
*── 「 TIKTOK 」 ──*

======================
*Username*: [username]
*Description*: [description]
======================
`
	SNAPINSTA_LIST = `
*── 「 INSTAGRAM 」 ──*

======================
*Username*: [username]
======================

Available Slide:

[slides]
======================
Balas nomor sebuah slide untuk mendownload, contoh:
*1* untuk memilih slide pertama
======================`

	SNAPINSTA_RESULT = `
*── 「 INSTAGRAM 」 ──*

======================
*Username*: [username]
======================
`
	SNAPTWIT_RESULT = `
*── 「 TWITTER 」 ──*

======================
*Username*: [username]
*Description*: [description]
======================
`

	SNAPSAVE_RESULT = `
*── 「 FACEBOOK 」 ──*

Available Quality:

[quality]
======================
Balas nomor urutan quality untuk mendownload, contoh:
*1* untuk memilih quality pertama
======================
`
)
