package models

import "time"

type ReportToday struct {
	TotalRevenue   int                  `json:"total_revenue"`
	TotalTransaksi int                  `json:"total_transaksi"`
	ProdukTerlaris []BestSellingProduct `json:"produk_terlaris"`
}

type ReportTodayCollect struct {
	Total     int       `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}
