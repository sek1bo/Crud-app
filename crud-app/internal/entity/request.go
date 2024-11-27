package entity

type Request struct {
	ID      int    `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Status  string `db:"status"`
}
