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
)
