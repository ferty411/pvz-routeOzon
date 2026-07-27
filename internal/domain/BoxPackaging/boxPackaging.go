package boxpackaging

import "pvz/internal/domain"

type BoxPackaging struct {
	Wrapper domain.OrderPackeger
}

func (b *BoxPackaging) GetPrice() float64 {
	return b.GetPrice() + domain.BoxPrice
}

func (b *BoxPackaging) ValidateWeight(weight float64) error {
	if weight >= domain.BoxWeightLimit {
		return domain.ErrWeightTooHeavy
	}
	return b.Wrapper.ValidateWeight(weight)
}
