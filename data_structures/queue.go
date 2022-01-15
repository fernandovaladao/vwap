package data_structures

//go:generate mockgen -destination queue_mock.go -package data_structures github.com/fernandovaladao/vwap/data_structures Queue
type Queue interface {
	Enqueue(elem float64) error
	Dequeue() (float64, error)
	IsFull() bool
	IsEmpty() bool
}
