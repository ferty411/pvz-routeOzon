package filmpackaging

import "pvz/internal/domain"

type FilmPackaging struct {
	Wrapper domain.OrderPackeger
}

func (f *FilmPackaging) GetPrice() float64 {
	return f.Wrapper.GetPrice() + domain.FilmPrice
}

func (f *FilmPackaging) ValidateWeight(weight float64) error {
	return f.Wrapper.ValidateWeight(weight)
}
