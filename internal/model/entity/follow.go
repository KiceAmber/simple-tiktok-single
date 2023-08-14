package entity

type Follow struct {
}

func (*Follow) TableName() string {
	return "follow"
}
