// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gexecutor_test

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/util/gconv"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gexecutor"
)

func TestExecutor_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test basic execution with input and output
		ctx := context.Background()
		executor := gexecutor.New[int, string](42)

		result, err := executor.WithMain(func(ctx context.Context, input int) (string, error) {
			return "result_" + gconv.String(input), nil
		}).Do(ctx)

		t.Assert(err, nil)
		t.Assert(result, "result_42")
	})
}

func TestExecutor_WithBefore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		executed := gtype.NewBool(false)
		ctx := context.Background()

		executor := gexecutor.New[int, string](10)
		executor = executor.WithBefore(func(ctx context.Context, input int) {
			executed.Set(true)
		}).WithMain(func(ctx context.Context, input int) (string, error) {
			return "processed", nil
		})

		result, err := executor.Do(ctx)
		t.Assert(err, nil)
		t.Assert(result, "processed")
		t.Assert(executed.Val(), true)
	})
}

func TestExecutor_WithAfter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		afterResult := gtype.NewString("")
		ctx := context.Background()

		executor := gexecutor.New[string, int]("test_input")
		executor = executor.WithMain(func(ctx context.Context, input string) (int, error) {
			return len(input), nil
		}).WithAfter(func(ctx context.Context, result int) {
			afterResult.Set("length_is_" + gconv.String(result))
		})

		result, err := executor.Do(ctx)

		t.Assert(err, nil)
		t.Assert(result, 10) // length of "test_input"
		t.Assert(afterResult.Val(), "length_is_10")
	})
}

func TestExecutor_WithBeforeAndAfter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		beforeExecuted := gtype.NewBool(false)
		afterResult := gtype.NewBool(false)
		ctx := context.Background()

		executor := gexecutor.New[float64, float64](5.5)
		executor = executor.WithBefore(func(ctx context.Context, input float64) {
			beforeExecuted.Set(true)
		}).WithMain(func(ctx context.Context, input float64) (float64, error) {
			return input * 2, nil
		}).WithAfter(func(ctx context.Context, result float64) {
			afterResult.Set(result == 11.0)
		})

		result, err := executor.Do(ctx)

		t.Assert(err, nil)
		t.Assert(result, 11.0)
		t.Assert(beforeExecuted.Val(), true)
		t.Assert(afterResult.Val(), true)
	})
}

func TestExecutor_MainFunctionError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expectedErr := errors.New("main function error")

		executor := gexecutor.New[string, string]("input")
		executor = executor.WithMain(func(ctx context.Context, input string) (string, error) {
			return "", expectedErr
		})

		result, err := executor.Do(context.Background())

		t.Assert(err, expectedErr)
		t.Assert(result, "")
	})
}

func TestExecutor_MainFunctionNotSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		executor := gexecutor.New[string, string]("input")

		result, err := executor.Do(context.Background())

		t.AssertNE(err, nil)
		t.Assert(result, "")
		t.Assert(err, gexecutor.ErrMainFuncNotSet)
	})
}

func TestExecutor_WithContext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with context that has values
		ctx := context.Background()
		ctx = context.WithValue(ctx, "key", "value")

		executor := gexecutor.New[int, string](100)
		executor = executor.WithMain(func(ctx context.Context, input int) (string, error) {
			ctxValue := ctx.Value("key").(string)
			res := "input_was_" + gconv.String(input) + "_ctx_value_" + ctxValue
			return res, nil
		})

		result, err := executor.Do(ctx)

		t.Assert(err, nil)
		t.Assert(result, "input_was_100_ctx_value_value")
	})
}

func TestExecutor_NilContext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		executor := gexecutor.New[string, int]("hello")
		executor = executor.WithMain(func(ctx context.Context, input string) (int, error) {
			return len(input), nil
		})

		result, err := executor.Do(nil) // Passing nil context should use background context

		t.Assert(err, nil)
		t.Assert(result, 5) // length of "hello"
	})
}

func TestExecutor_ChainUsage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		beforeCalled := gtype.NewBool(false)
		afterCalled := gtype.NewBool(false)

		result, err := gexecutor.New[int, string](42).
			WithBefore(func(ctx context.Context, input int) {
				beforeCalled.Set(true)
			}).
			WithMain(func(ctx context.Context, input int) (string, error) {
				return "value_" + gconv.String(input), nil
			}).
			WithAfter(func(ctx context.Context, result string) {
				afterCalled.Set(true)
			}).
			Do(context.Background())

		t.Assert(err, nil)
		t.Assert(result, "value_42")
		t.Assert(beforeCalled.Val(), true)
		t.Assert(afterCalled.Val(), true)
	})
}

func TestExecutor_GenericTypes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with complex types
		type InputStruct struct {
			ID   int
			Name string
		}

		type OutputStruct struct {
			ProcessedID   int
			ProcessedName string
			Status        bool
		}

		input := InputStruct{ID: 123, Name: "test"}

		executor := gexecutor.New[InputStruct, OutputStruct](input)
		executor = executor.WithBefore(func(ctx context.Context, input InputStruct) {
			// Pre-processing
		}).WithMain(func(ctx context.Context, input InputStruct) (OutputStruct, error) {
			return OutputStruct{
				ProcessedID:   input.ID * 2,
				ProcessedName: input.Name + "_processed",
				Status:        true,
			}, nil
		}).WithAfter(func(ctx context.Context, result OutputStruct) {
			// Post-processing
		})

		result, err := executor.Do(context.Background())

		t.Assert(err, nil)
		t.Assert(result.ProcessedID, 246)
		t.Assert(result.ProcessedName, "test_processed")
		t.Assert(result.Status, true)
	})
}

func TestExecutor_TemplateReuse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		// Test that WithXxx methods return new instances, allowing template reuse
		templateExecutor := gexecutor.New[int, string](10) // Use 10 as base input
		templateExecutor = templateExecutor.WithBefore(func(ctx context.Context, input int) {
			// Common before logic
		}).WithAfter(func(ctx context.Context, result string) {
			// Common after logic
		})

		// First specialized executor - override main function but reuse before/after
		executor1 := templateExecutor.WithMain(func(ctx context.Context, input int) (string, error) {
			return "doubled: " + gconv.String(input*2), nil
		})
		result1, err1 := executor1.Do(ctx)
		t.Assert(err1, nil)
		t.Assert(result1, "doubled: 20") // with input 10

		// Second specialized executor - reuse same before/after but different main function
		executor2 := templateExecutor.WithMain(func(ctx context.Context, input int) (string, error) {
			return "tripled: " + gconv.String(input*3), nil
		})
		result2, err2 := executor2.Do(ctx)
		t.Assert(err2, nil)
		t.Assert(result2, "tripled: 30") // with input 10

		// Third executor created from scratch with different input and reused patterns
		executor3 := gexecutor.New[int, string](5).WithBefore(func(ctx context.Context, input int) {
			// Same before logic as template
		}).
			WithAfter(func(ctx context.Context, result string) {
				// Same after logic as template
			}).
			WithMain(func(ctx context.Context, input int) (string, error) {
				return "squared: " + gconv.String(input*input), nil
			})
		result3, err3 := executor3.Do(ctx)
		t.Assert(err3, nil)
		t.Assert(result3, "squared: 25") // with input 5

		// Show that original template executor still works with its own input
		executor4 := templateExecutor.WithMain(func(ctx context.Context, input int) (string, error) {
			return "original: " + gconv.String(input), nil
		})
		result4, err4 := executor4.Do(ctx)
		t.Assert(err4, nil)
		t.Assert(result4, "original: 10") // with original input 10
	})
}
