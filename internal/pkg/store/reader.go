package store

// TODO: Create Reader struct/interface to abstract the file reading process
type Reader struct{}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) Scan(f func(v interface{}, err error) error) error {
	// Do something...
	_ = f(nil, nil)
	return nil
}
