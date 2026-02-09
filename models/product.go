package models

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Category_ID int    `json:"category_id"`
}

type BestSellingProduct struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	QtyTerjual int    `json:"qty_terjual"`
}
