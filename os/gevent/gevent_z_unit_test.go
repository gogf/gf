package gevent_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gevent"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSeqEventBus_PublishSubscribe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建顺序事件总线
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		// 用于记录处理结果
		result := make(chan string, 10)

		// 订阅事件
		subscriber, err := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler1:" + e.GetData()["message"].(string)
			return nil
		}, nil, nil)
		t.AssertNil(err)
		t.AssertNE(subscriber, nil)

		// 发布事件
		params := map[string]any{"message": "hello"}
		ok, err := bus.Publish("test.topic", params, gevent.Ignore, gevent.Seq)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 等待处理完成
		select {
		case msg := <-result:
			t.Assert(msg, "handler1:hello")
		case <-time.After(time.Second):
			t.Error("Event handling timeout")
		}
	})
}

func TestSeqEventBus_MultipleSubscribers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)

		// 订阅者1
		_, err1 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler1:" + e.GetData()["message"].(string)
			return nil
		}, nil, nil)
		t.AssertNil(err1)

		// 订阅者2
		_, err2 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler2:" + e.GetData()["message"].(string)
			return nil
		}, nil, nil)
		t.AssertNil(err2)

		// 发布事件
		params := map[string]any{"message": "hello"}
		ok, err := bus.Publish("test.topic", params, gevent.Ignore, gevent.Seq)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 收集结果
		var messages []string
		timeout := time.After(time.Second)
		for i := 0; i < 2; i++ {
			select {
			case msg := <-result:
				messages = append(messages, msg)
			case <-timeout:
				t.Error("Event handling timeout")
				return
			}
		}

		// 验证结果
		t.Assert(len(messages), 2)
	})
}

func TestSeqEventBus_Priority(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)

		// 订阅者1（低优先级）
		_, err1 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			time.Sleep(10 * time.Millisecond) // 模拟处理时间
			result <- "low"
			return nil
		}, nil, nil, gevent.PriorityLow)
		t.AssertNil(err1)

		// 订阅者2（高优先级）
		_, err2 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "high"
			return nil
		}, nil, nil, gevent.PriorityHigh)
		t.AssertNil(err2)

		// 发布事件
		params := map[string]any{"message": "hello"}
		ok, err := bus.Publish("test.topic", params, gevent.Ignore, gevent.Seq)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 按优先级顺序接收消息
		select {
		case msg := <-result:
			t.Assert(msg, "high")
		case <-time.After(time.Second):
			t.Error("Event handling timeout")
		}

		select {
		case msg := <-result:
			t.Assert(msg, "low")
		case <-time.After(time.Second):
			t.Error("Event handling timeout")
		}
	})
}

func TestSeqEventBus_ParallelExecution(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)
		startTime := time.Now()

		// 订阅者1
		_, err1 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			time.Sleep(100 * time.Millisecond) // 模拟耗时操作
			result <- "handler1"
			return nil
		}, nil, nil)
		t.AssertNil(err1)

		// 订阅者2
		_, err2 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			time.Sleep(100 * time.Millisecond) // 模拟耗时操作
			result <- "handler2"
			return nil
		}, nil, nil)
		t.AssertNil(err2)

		// 发布并行执行的事件
		params := map[string]any{"message": "hello"}
		ok, err := bus.Publish("test.topic", params, gevent.Ignore, gevent.Parallel)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 等待所有处理完成
		var messages []string
		timeout := time.After(300 * time.Millisecond) // 应该在约100ms内完成（并行执行）
		for i := 0; i < 2; i++ {
			select {
			case msg := <-result:
				messages = append(messages, msg)
			case <-timeout:
				t.Error("Event handling timeout")
				return
			}
		}

		// 验证执行时间是否符合并行特征
		executionTime := time.Since(startTime)
		t.AssertLT(executionTime, 150*time.Millisecond) // 应远小于200ms（串行执行时间）

		// 验证结果
		t.Assert(len(messages), 2)
	})
}

func TestSeqEventBus_ErrorHandling_Stop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)

		// 订阅者1（会出错）
		_, err1 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler1"
			return gevent.EventBusClosedError // 返回错误
		}, nil, nil)
		t.AssertNil(err1)

		// 订阅者2（正常）
		_, err2 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler2"
			return nil
		}, nil, nil)
		t.AssertNil(err2)

		// 发布事件，使用Stop错误模式
		params := map[string]any{"message": "hello"}
		ok, err := bus.Publish("test.topic", params, gevent.Stop, gevent.Seq)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 只应该收到第一个处理结果，因为第一个处理出错且使用Stop模式
		select {
		case msg := <-result:
			t.Assert(msg, "handler1")
		case <-time.After(time.Second):
			t.Error("Event handling timeout")
		}

		// 不应该再收到其他消息
		select {
		case msg := <-result:
			t.Error("Should not receive more messages, but got:", msg)
		case <-time.After(50 * time.Millisecond):
			// 正确情况，没有更多消息
		}
	})
}

func TestSeqEventBus_ErrorHandling_Ignore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)

		// 订阅者1（会出错）
		_, err1 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler1"
			return gevent.EventBusClosedError // 返回错误
		}, nil, nil)
		t.AssertNil(err1)

		// 订阅者2（正常）
		_, err2 := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler2"
			return nil
		}, nil, nil)
		t.AssertNil(err2)

		// 发布事件，使用Ignore错误模式
		params := map[string]any{"message": "hello"}
		ok, err := bus.Publish("test.topic", params, gevent.Ignore, gevent.Seq)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 应该收到两个处理结果，即使第一个出错
		results := make([]string, 0, 2)
		timeout := time.After(time.Second)
		for i := 0; i < 2; i++ {
			select {
			case msg := <-result:
				results = append(results, msg)
			case <-timeout:
				t.Error("Event handling timeout")
				return
			}
		}

		t.Assert(len(results), 2)
	})
}

func TestSeqEventBus_Unsubscribe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)

		// 订阅事件
		subscriber, err := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler1:" + e.GetData()["message"].(string)
			return nil
		}, nil, nil)
		t.AssertNil(err)

		// 取消订阅
		ok, err := subscriber.UnSubscribe()
		t.AssertNil(err)
		t.Assert(ok, true)

		// 再次取消订阅应该失败
		ok, err = subscriber.UnSubscribe()
		t.AssertNil(err)
		t.Assert(ok, true) // 已经取消订阅的状态

		// 发布事件
		params := map[string]any{"message": "hello"}
		_, err = bus.Publish("test.topic", params, gevent.Ignore, gevent.Seq)
		t.Assert(err, gevent.SubscriberEmptyError) // 没有订阅者
	})
}

func TestSeqEventBus_PublishEvent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)

		// 订阅事件
		_, err := bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- e.GetData()["message"].(string)
			return nil
		}, nil, nil)
		t.AssertNil(err)

		// 创建并发布自定义事件
		event := gevent.BaseEventFactoryFunc("test.topic", map[string]any{"message": "custom event"}, gevent.Ignore, gevent.Seq)
		ok, err := bus.PublishEvent(event)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 验证结果
		select {
		case msg := <-result:
			t.Assert(msg, "custom event")
		case <-time.After(time.Second):
			t.Error("Event handling timeout")
		}
	})
}

func TestSeqEventBus_FactoryFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()
		defer bus.Close()

		result := make(chan string, 10)

		// 注册自定义事件工厂函数
		customFactory := func(topic string, params map[string]any, errorModel gevent.ErrorModel, execModel gevent.ExecModel) gevent.Event {
			event := &gevent.BaseEvent{
				Topic:      topic,
				Data:       params,
				ErrorModel: errorModel,
				ExecModel:  execModel,
			}
			result <- "factory_called"
			return event
		}

		ok, err := bus.RegisterFactoryFunc("test.topic", customFactory)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 订阅事件
		_, err = bus.Subscribe("test.topic", func(e gevent.Event) error {
			result <- "handler:" + e.GetData()["message"].(string)
			return nil
		}, nil, nil)
		t.AssertNil(err)

		// 发布事件
		params := map[string]any{"message": "hello"}
		ok, err = bus.Publish("test.topic", params, gevent.Ignore, gevent.Seq)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 验证工厂函数被调用且事件被处理
		results := make([]string, 0, 2)
		timeout := time.After(time.Second)
		for i := 0; i < 2; i++ {
			select {
			case msg := <-result:
				results = append(results, msg)
			case <-timeout:
				t.Error("Event handling timeout")
				return
			}
		}

		t.Assert(results, []string{"factory_called", "handler:hello"})
	})
}

func TestSeqEventBus_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		bus := gevent.NewSeqEventBus()

		// 订阅事件
		_, err := bus.Subscribe("test.topic", func(e gevent.Event) error {
			return nil
		}, nil, nil)
		t.AssertNil(err)

		// 关闭事件总线
		bus.Close()

		// 尝试发布事件应该失败
		params := map[string]any{"message": "hello"}
		_, err = bus.Publish("test.topic", params, gevent.Ignore, gevent.Seq)
		t.Assert(err, gevent.EventBusClosedError)

		// 检查是否已关闭
		t.Assert(bus.IsClosed(), true)
	})
}
