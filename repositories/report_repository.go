package repositories

import (
	"category-api/models"
	"database/sql"
	"errors"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReportToday() (*models.ReportToday, error) {
	// get total revenue
	now := time.Now()

	startOfToday := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(), // penting: samakan dengan timezone DB
	)

	startOfTomorrow := startOfToday.AddDate(0, 0, 1)

	query := "SELECT COALESCE(SUM(total_amount), 0) as total FROM transactions WHERE created_at >= $1 AND created_at < $2"
	var totalRevenue int
	err := repo.db.QueryRow(query, startOfToday, startOfTomorrow).Scan(&totalRevenue)
	if err == sql.ErrNoRows {
		return nil, errors.New("no sales today")
	}
	if totalRevenue == 0 {
		return nil, errors.New("no sales today")
	}
	if err != nil {
		return nil, err
	}

	query_2 := "SELECT SUM(x.row) FROM(SELECT 1 as row FROM transactions WHERE created_at >= $1 AND created_at < $2) as x"
	var totalTransaksi int
	err_2 := repo.db.QueryRow(query_2, startOfToday, startOfTomorrow).Scan(&totalTransaksi)
	if err_2 == sql.ErrNoRows {
		return nil, errors.New("no sales today")
	}
	if err_2 != nil {
		return nil, err_2
	}

	query_3 := "SELECT pr.id, pr.name, SUM(td.quantity) as qty FROM transactions as t INNER JOIN transaction_details td ON td.transaction_id = t.id INNER JOIN products pr ON pr.id = td.product_id WHERE t.created_at >= $1 AND t.created_at < $2 GROUP BY pr.id,pr.name ORDER BY qty DESC LIMIT 1"
	var produkTerlaris []models.BestSellingProduct
	var p models.BestSellingProduct
	err_3 := repo.db.QueryRow(query_3, startOfToday, startOfTomorrow).Scan(&p.ID, &p.Name, &p.QtyTerjual)
	if err_3 == sql.ErrNoRows {
		return nil, errors.New("no sales today")
	}
	if err_3 != nil {
		return nil, err_3
	}

	produkTerlaris = append(produkTerlaris, p)

	return &models.ReportToday{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: produkTerlaris,
	}, nil
}
