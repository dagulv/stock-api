package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

type StreamEncoder interface {
	EncodeToStream(stream *jsoniter.Stream)
}

type DomainEncoder[Model StreamEncoder] struct {
	Stream *jsoniter.Stream
	count  int
	// total  int TODO
}

func CreateDomainEncoder[Model StreamEncoder](json jsoniter.API, writer io.Writer) *DomainEncoder[Model] {
	stream := json.BorrowStream(writer)

	stream.WriteObjectStart()
	stream.WriteObjectField("result")
	stream.WriteArrayStart()

	return &DomainEncoder[Model]{
		Stream: stream,
	}
}

func (e *DomainEncoder[Model]) AddLine(model Model) {
	if e.count > 0 {
		e.Stream.WriteMore()
	}

	e.Stream.WriteObjectStart()

	model.EncodeToStream(e.Stream)

	e.Stream.WriteObjectEnd()

	e.count++
}

func (e *DomainEncoder[Model]) Flush() error {
	e.Stream.WriteArrayEnd()
	e.Stream.WriteMore()
	e.Stream.WriteObjectField("count")
	e.Stream.WriteInt(e.count)
	e.Stream.WriteObjectEnd()

	return e.Stream.Flush()
}
