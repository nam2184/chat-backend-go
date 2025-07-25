package response

type Meta struct {
  Total    *int    `json:"total"`
}

func NewMeta() *Meta {
	return &Meta{}
}

func (m *Meta) WithTotal(total int) {
    m.Total = &total
}




