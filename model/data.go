package model

type Data struct {
	Tanggal    string `json:"tanggal"`
	NamaVendor string `json:"nama_vendor"`
	KodeBarang string `json:"kode_barang"`
	CaraBayar  string `json:"cara_bayar"`
	QTY        int    `json:"qty"`
	Harga      int    `json:"harga"`
	Jumlah     int    `json:"jumlah"`
	Flag       string `json:"flag"`
}
