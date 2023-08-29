package common

import "errors"

var ErrWrongMetricsType = errors.New("wrong metrics type")
var ErrWrongMetricsName = errors.New("wrong metrics name")
var ErrWrongMetricsValue = errors.New("wrong metrics value")
