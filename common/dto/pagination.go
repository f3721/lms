package dto

type Pagination struct {
	PageIndex int `form:"pageIndex"`
	PageSize  int `form:"pageSize"`
}

func (m *Pagination) GetPageIndex() int {
	if m.PageIndex <= 0 {
		m.PageIndex = 1
	}
	return m.PageIndex
}

func (m *Pagination) GetPageSize() int {
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}

func (m *Pagination) GetPageSizeNegative() int {
	if m.PageSize < 0 {
		m.PageSize = -1
	} else if m.PageSize == 0 {
		m.PageSize = 10
	}
	return m.PageSize
}
