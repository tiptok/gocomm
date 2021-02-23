package trace

import (
	"context"
	"fmt"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/trace/tracespec"
	"github.com/tiptok/gocomm/xstring"
	"hash/crc32"
	"strconv"
	"strings"
	"time"
	//"github.com/tal-tech/go-zero/core/xstring"
	//"github.com/tal-tech/go-zero/core/timex"
)

const (
	initSpanId  = "0"
	clientFlag  = "client"
	serverFlag  = "server"
	spanSepRune = '.'
	timeFormat  = "2006-01-02 15:04:05.000"
)

var spanSep = string([]byte{spanSepRune})

type Span struct {
	ctx           spanContext
	serviceName   string
	operationName string
	startTime     time.Time
	flag          string
	children      int
}

func newServerSpan(carrier Carrier, serviceName, operationName string) tracespec.Trace {
	traceId := xstring.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(traceIdKey)
		}
		return ""
	}, func() string {
		return xstring.RandId()
	})
	spanId := xstring.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(spanIdKey)
		}
		return ""
	}, func() string {
		return initSpanId
	})

	return &Span{
		ctx: spanContext{
			traceId: traceId,
			spanId:  spanId,
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     time.Now(),
		flag:          serverFlag,
	}
}

func (s *Span) SpanId() string {
	return s.ctx.SpanId()
}

func (s *Span) TraceId() string {
	return s.ctx.TraceId()
}

func (s *Span) Visit(fn func(key, val string) bool) {
	s.ctx.Visit(fn)
}

func (s *Span) Finish() {
	if globalReport != nil {
		globalSpanFinish(s)
	}
}

func (s *Span) Fork(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId:  s.forkSpanId(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     time.Now(),
		flag:          clientFlag,
	}
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}

func (s *Span) Follow(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId:  s.followSpanId(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     time.Now(),
		flag:          s.flag,
	}
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}

func (s *Span) forkSpanId() string {
	s.children++
	return fmt.Sprintf("%s.%d", s.ctx.spanId, s.children)
}

func (s *Span) followSpanId() string {
	fields := strings.FieldsFunc(s.ctx.spanId, func(r rune) bool {
		return r == spanSepRune
	})
	if len(fields) == 0 {
		return s.ctx.spanId
	}

	last := fields[len(fields)-1]
	val, err := strconv.Atoi(last)
	if err != nil {
		return s.ctx.spanId
	}

	last = strconv.Itoa(val + 1)
	fields[len(fields)-1] = last

	return strings.Join(fields, spanSep)
}

/*******************自定义 tracing ******************/
func (s *Span) printInfo() {
	span := s
	log.Debug(span.TraceId(), span.SpanId(), span.serviceName, span.operationName, span.startTime.Unix())
}

func (s *Span) parentId() string {
	idx := strings.LastIndex(s.SpanId(), ".")
	if idx < 0 {
		return s.TraceId()
	}
	return string(s.SpanId()[:idx])
}

func zipkinOnFinish(s *Span) {
	parentId := model.ID(uint64(crc32.ChecksumIEEE([]byte(s.parentId()))))
	m := model.SpanModel{
		SpanContext: model.SpanContext{
			TraceID: model.TraceID{
				Low: uint64(crc32.ChecksumIEEE([]byte(s.TraceId()))),
			},
			ID: model.ID(uint64(crc32.ChecksumIEEE([]byte(s.SpanId())))),
		},
		Name:          s.operationName,
		Kind:          model.Kind(strings.ToUpper(s.flag)),
		Timestamp:     time.Now(),
		Duration:      time.Now().Sub(s.startTime),
		LocalEndpoint: globalLocalEndpoint,
		Annotations:   make([]model.Annotation, 0),
		Tags:          make(map[string]string),
	}
	if strings.LastIndex(s.SpanId(), ".") >= 0 {
		m.ParentID = &parentId
	}
	//s.printInfo()
	log.Error(common.JsonAssertString(m))
	globalReport.Send(m)
}
