package storage

import (
	"encoding/json"
	"os"
	"pvz/internal/domain"
)

type FileStorage struct {
	filename string
	orders   map[string]domain.Order
}

func NewFileStorage(path string) (*FileStorage, error) {
	fs := &FileStorage{
		filename: path,
		orders:   make(map[string]domain.Order),
	}

	err := fs.LoadFromFile()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStorage) LoadFromFile() error {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		return err
	}

	if len(data) <= 0 {
		return nil
	}

	return json.Unmarshal(data, &fs.orders)
}
func (fs *FileStorage) SaveToFile() error {
	data, err := json.MarshalIndent(fs.orders, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filename, data, 0644)
}

func (fs *FileStorage) GetOrder(orderID string) (domain.Order, error) {
	order, ok := fs.orders[orderID]
	if !ok {
		return domain.Order{}, domain.ErrOrderNotFound
	}
	return order, nil
}

func (fs *FileStorage) DeleteOrder(orderID string) error {
	if _, ok := fs.orders[orderID]; !ok {
		return domain.ErrOrderNotFound
	}
	delete(fs.orders, orderID)
	return nil
}

func (fs *FileStorage) SaveOrder(order domain.Order) error {
	fs.orders[order.OrderID] = order
	return fs.SaveToFile()
}

func (fs *FileStorage) GetAllOrders() ([]domain.Order, error) {
	var list []domain.Order

	for _, order := range fs.orders {
		list = append(list, order)
	}
	return list, nil
}
