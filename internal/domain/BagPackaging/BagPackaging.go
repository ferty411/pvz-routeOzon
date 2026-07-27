package bagpackaging

import "pvz/internal/domain"

type BagPackaging struct {
	Wrapper domain.OrderPackeger
}

func (b *BagPackaging) GetPrice() float64 {
	return b.GetPrice() + domain.BagPrice
}

func (b *BagPackaging) ValidateWeight(weight float64) error {
	if weight >= domain.BagWeightLimit {
		return domain.ErrWeightTooHeavy
	}
	return b.Wrapper.ValidateWeight(weight)
}
