package gevent_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gevent"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestEvent_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()
		t.AssertNE(event, nil)
	})
}

func TestEvent_SubscribeAndPublish(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		var result string
		handler := func(topic string, message any) {
			result = message.(string)
		}

		sub, err := event.Subscribe("test", handler)
		t.AssertNil(err)
		t.AssertNE(sub, nil)

		err = event.Publish("test", "Hello World")
		t.AssertNil(err)

		time.Sleep(10 * time.Millisecond)
		t.Assert(result, "Hello World")
	})
}

func TestEvent_SubscribeSync(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		var result string
		handler := func(topic string, message any) {
			result = message.(string)
		}

		sub, err := event.Subscribe("test", handler)
		t.AssertNil(err)
		t.AssertNE(sub, nil)

		err = event.PublishSync("test", "Hello Sync")
		t.AssertNil(err)

		t.Assert(result, "Hello Sync")
	})
}

func TestEvent_Priority(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		var results []string
		handlerLow := func(topic string, message any) {
			results = append(results, "low")
		}
		handlerNormal := func(topic string, message any) {
			results = append(results, "normal")
		}
		handlerHigh := func(topic string, message any) {
			results = append(results, "high")
		}
		handlerUrgent := func(topic string, message any) {
			results = append(results, "urgent")
		}
		handlerImmediate := func(topic string, message any) {
			results = append(results, "immediate")
		}

		_, err := event.Subscribe("test", handlerLow, gevent.PriorityLow)
		t.AssertNil(err)

		_, err = event.Subscribe("test", handlerNormal, gevent.PriorityNormal)
		t.AssertNil(err)

		_, err = event.Subscribe("test", handlerHigh, gevent.PriorityHigh)
		t.AssertNil(err)

		_, err = event.Subscribe("test", handlerUrgent, gevent.PriorityUrgent)
		t.AssertNil(err)

		_, err = event.Subscribe("test", handlerImmediate, gevent.PriorityImmediate)
		t.AssertNil(err)

		err = event.PublishSync("test", "message")
		t.AssertNil(err)

		t.Assert(results, []string{"immediate", "urgent", "high", "normal", "low"})
	})
}

// 添加一个新的测试用例，测试相同优先级下的执行顺序
func TestEvent_PrioritySameLevel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		var results []int
		handler1 := func(topic string, message any) {
			results = append(results, 1)
		}
		handler2 := func(topic string, message any) {
			results = append(results, 2)
		}
		handler3 := func(topic string, message any) {
			results = append(results, 3)
		}

		// 按顺序订阅相同优先级的处理器
		_, err := event.Subscribe("test", handler1, gevent.PriorityHigh)
		t.AssertNil(err)

		_, err = event.Subscribe("test", handler2, gevent.PriorityHigh)
		t.AssertNil(err)

		_, err = event.Subscribe("test", handler3, gevent.PriorityHigh)
		t.AssertNil(err)

		err = event.PublishSync("test", "message")
		t.AssertNil(err)

		t.Assert(results, []int{1, 2, 3})
	})
}

func TestEvent_Unsubscribe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		var count int
		handler := func(topic string, message any) {
			count++
		}

		sub, err := event.Subscribe("test", handler)
		t.AssertNil(err)

		err = event.PublishSync("test", "message")
		t.AssertNil(err)
		t.Assert(count, 1)

		sub.Unsubscribe()

		err = event.PublishSync("test", "message")
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func TestEvent_UnsubscribeMultipleTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		var count int
		handler := func(topic string, message any) {
			count++
		}

		sub, err := event.Subscribe("test", handler)
		t.AssertNil(err)

		err = event.PublishSync("test", "message")
		t.AssertNil(err)
		t.Assert(count, 1)

		sub.Unsubscribe()
		sub.Unsubscribe()
		sub.Unsubscribe()

		err = event.PublishSync("test", "message")
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func TestEvent_SubscribersCount(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		handler := func(topic string, message any) {}

		t.Assert(event.SubscribersCount("test"), 0)

		sub1, err := event.Subscribe("test", handler)
		t.AssertNil(err)
		t.Assert(event.SubscribersCount("test"), 1)

		sub2, err := event.Subscribe("test", handler)
		t.AssertNil(err)
		t.Assert(event.SubscribersCount("test"), 2)

		sub1.Unsubscribe()
		t.Assert(event.SubscribersCount("test"), 1)

		sub2.Unsubscribe()
		t.Assert(event.SubscribersCount("test"), 0)
	})
}

func TestEvent_Topics(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		handler := func(topic string, message any) {}

		t.Assert(len(event.Topics()), 0)

		_, err := event.Subscribe("test1", handler)
		t.AssertNil(err)

		_, err = event.Subscribe("test2", handler)
		t.AssertNil(err)

		topics := event.Topics()
		t.Assert(len(topics), 2)
		t.AssertIN("test1", topics)
		t.AssertIN("test2", topics)
	})
}

func TestEvent_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		var count int
		handler := func(topic string, message any) {
			count++
		}

		_, err := event.Subscribe("test", handler)
		t.AssertNil(err)

		err = event.PublishSync("test", "message")
		t.AssertNil(err)
		t.Assert(count, 1)

		event.Clear()

		err = event.PublishSync("test", "message")
		t.AssertNil(err)
		t.Assert(count, 1)
		t.Assert(event.SubscribersCount("test"), 0)
	})
}

func TestEvent_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		event.Close()

		handler := func(topic string, message any) {}

		_, err := event.Subscribe("test", handler)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "event manager is closed")

		err = event.Publish("test", "message")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "event manager is closed")
	})
}

func TestEvent_RecoverFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		panicHandled := gtype.NewBool()
		recoverFunc := func(id int64, topic string, message any, err any) {
			panicHandled.Set(true)
		}

		handler := func(topic string, message any) {
			panic("test panic")
		}

		_, err := event.SubscribeWithRecover("test", handler, recoverFunc)
		t.AssertNil(err)

		err = event.PublishSync("test", "message")
		t.AssertNil(err)
		t.Assert(panicHandled, true)
	})
}

func TestEvent_RecoverFuncAsync(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		panicHandled := gtype.NewBool()
		recoverFunc := func(id int64, topic string, message any, err any) {
			panicHandled.Set(true)
		}

		handler := func(topic string, message any) {
			panic("test panic")
		}

		_, err := event.SubscribeWithRecover("test", handler, recoverFunc)
		t.AssertNil(err)

		err = event.Publish("test", "message")
		t.AssertNil(err)

		time.Sleep(10 * time.Millisecond)
		t.Assert(panicHandled.Val(), true)
	})
}

func TestEvent_SubscribeAndPublishWithoutTimeDependency(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		resultChan := make(chan string, 1)
		handler := func(topic string, message any) {
			resultChan <- message.(string)
		}

		sub, err := event.Subscribe("test", handler)
		t.AssertNil(err)
		t.AssertNE(sub, nil)

		err = event.Publish("test", "Hello World")
		t.AssertNil(err)

		// 使用通道等待异步处理完成，而不是固定的时间等待
		select {
		case result := <-resultChan:
			t.Assert(result, "Hello World")
		case <-time.After(1 * time.Second):
			t.Error("Event handler was not called within timeout")
		}
	})
}
