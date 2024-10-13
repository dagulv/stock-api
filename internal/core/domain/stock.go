package domain

import (
	jsoniter "github.com/json-iterator/go"
)

type Stock struct {
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float32 `json:"int"`
}

func (st *Stock) EncodeToStream(s *jsoniter.Stream) {
	s.WriteObjectField("symbol")
	s.WriteString(st.Symbol)

	s.WriteMore()
	s.WriteObjectField("name")
	s.WriteString(st.Name)

	s.WriteMore()
	s.WriteObjectField("price")
	s.WriteFloat32(st.Price)
}
