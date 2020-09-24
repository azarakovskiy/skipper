package tracing

import (
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/zalando/skipper/filters"
	logfilter "github.com/zalando/skipper/filters/log"
)

const (
	StateBagToTagFilterName = "stateBagToTag"
)

type stateBagToTagSpec struct{}

type stateBagToTagFilter struct {
	stateBagItemName string
	tagName          string
}

func (stateBagToTagSpec) Name() string {
	return StateBagToTagFilterName
}

func (stateBagToTagSpec) CreateFilter(args []interface{}) (filters.Filter, error) {
	if len(args) < 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	stateBagItemName, ok := args[0].(string)
	if !ok || stateBagItemName == "" {
		return nil, filters.ErrInvalidFilterParameters
	}

	tagName := stateBagItemName
	if len(args) > 1 {
		tagNameArg, ok := args[1].(string)
		if !ok || tagNameArg == "" {
			return nil, filters.ErrInvalidFilterParameters
		}
		tagName = tagNameArg
	}

	return stateBagToTagFilter{
		stateBagItemName: stateBagItemName,
		tagName:          tagName,
	}, nil
}

func NewStateBagToTag() filters.Spec {
	return stateBagToTagSpec{}
}

func (f stateBagToTagFilter) Request(ctx filters.FilterContext) {
	span := opentracing.SpanFromContext(ctx.Request().Context())
	if span == nil {
		return
	}

	if f.stateBagItemName == logfilter.AuthUserKey {
		value, ok := ctx.StateBag()[logfilter.MaskedAuthUserKey]
		if ok {
			span.SetTag(f.tagName, value.(string))
			return
		}
	}

	value, ok := ctx.StateBag()[f.stateBagItemName]
	if !ok {
		return
	}
	span.SetTag(f.tagName, fmt.Sprint(value))
}

func (stateBagToTagFilter) Response(ctx filters.FilterContext) {}
