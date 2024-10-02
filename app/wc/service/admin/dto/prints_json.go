package dto

type CommonPrintsReq struct {
	Id        int    `uri:"id"`
	PrintType string `uri:"print-type"`
}

func (s *CommonPrintsReq) GetId() interface{} {
	return s.Id
}

type CommonPrintsResp struct {
	Prints any `json:"prints"`
}
