package echor

// Query struct
type Query struct {
	Where  string
	Offset int64
	Limit  int64
	Order  string
}

// InsertModelOpts struct
type InsertModelOpts struct {
	skipIfConflict bool
}
