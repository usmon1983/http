package banners

import (
	"errors"
	"sync"
	"context"
)

//Service представляет собой сервис по управлению баннерами.
type Service struct {
	mu sync.RWMutex
	items []*Banner
}

//NewService создаёт сервис.
func NewService() *Service  {
	return &Service{items: make([]*Banner, 0)}
}

//Banner представляет собой баннер.
type Banner struct {
	ID int64
	Title string
	Content string
	Button string
	Link string
}

//All возвращает все существующие баннеры.
func (s *Service) All(ctx context.Context) ([]*Banner, error)  {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

//ByID возвращает баннер по идентификатору.
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error)  {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.items {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, errors.New("item not found")
}

var BannerID int64 = 0
//Save сохраняет/обновляет баннер.
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error)  {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item.ID == 0 {
		BannerID++
		item.ID = BannerID
		s.items = append(s.items, item)
		return item, nil
	}
	for k, v := range s.items {
		if v.ID == item.ID {
			s.items[k] = item
			return item, nil
		}
	}
	return nil, errors.New("item not found")
}

//RemoveByID удаляет баннер по идентификатору.
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error)  {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for k, v := range s.items {
		if v.ID == id {
			s.items = append(s.items[:k], s.items[k+1:]...)
			return v, nil
		}
	}

	return nil, errors.New("item not found")
}