package baseorder

type BaseOrder struct {
	InitialPrice float64
}

func (b *BaseOrder) GetPrice() float64                   { return b.InitialPrice }
func (b *BaseOrder) ValidateWeight(weight float64) error { return nil }
