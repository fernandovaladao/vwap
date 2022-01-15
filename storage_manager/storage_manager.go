package storage_manager

//go:generate mockgen -destination storage_manager_mock.go -package storage_manager github.com/fernandovaladao/vwap/storage_manager StorageManager
type StorageManager interface {
	Store(price float64) error
	GetSum() float64
}
