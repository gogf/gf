package gevent_test

import (
	"testing"
	"time"

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
		received := make(chan string, 1)

		// Subscribe to a topic
		subscriber := event.Subscribe("test.topic", func(topic string, message any) {
			received <- message.(string)
		})

		// Publish a message
		event.Publish("test.topic", "Hello World")

		// Wait for the message to be received
		select {
		case msg := <-received:
			t.Assert(msg, "Hello World")
		case <-time.After(time.Second):
			t.Error("Message not received within timeout")
		}

		// Test unsubscribe
		subscriber.Unsubscribe()
		t.Assert(event.SubscribersCount("test.topic"), 0)
	})
}

func TestEvent_SubscribeFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()
		received := make(chan string, 1)

		// Subscribe using SubscribeFunc
		event.SubscribeFunc("test.topic", func(topic string, message any) {
			received <- message.(string)
		})

		// Publish a message
		event.Publish("test.topic", "Hello from SubscribeFunc")

		// Wait for the message to be received
		select {
		case msg := <-received:
			t.Assert(msg, "Hello from SubscribeFunc")
		case <-time.After(time.Second):
			t.Error("Message not received within timeout")
		}
	})
}

func TestEvent_PublishSync(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()
		executionOrder := make([]int, 0)

		// Subscribe multiple handlers
		event.Subscribe("sync.topic", func(topic string, message any) {
			executionOrder = append(executionOrder, 1)
		})

		event.Subscribe("sync.topic", func(topic string, message any) {
			executionOrder = append(executionOrder, 2)
		})

		// Publish synchronously
		event.PublishSync("sync.topic", "test message")

		// Check that execution order is deterministic
		t.Assert(len(executionOrder), 2)
		t.Assert(executionOrder[0], 1)
		t.Assert(executionOrder[1], 2)
	})
}

func TestEvent_Priority(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()
		executionOrder := make([]int, 0)

		// Subscribe handlers with different priorities
		event.Subscribe("priority.topic", func(topic string, message any) {
			executionOrder = append(executionOrder, 3) // Normal priority
		}, gevent.PriorityNormal)

		event.Subscribe("priority.topic", func(topic string, message any) {
			executionOrder = append(executionOrder, 1) // High priority
		}, gevent.PriorityHigh)

		event.Subscribe("priority.topic", func(topic string, message any) {
			executionOrder = append(executionOrder, 2) // Medium priority
		}, gevent.PriorityUrgent)

		event.Subscribe("priority.topic", func(topic string, message any) {
			executionOrder = append(executionOrder, 4) // Low priority
		}, gevent.PriorityLow)

		// Publish synchronously to ensure order
		event.PublishSync("priority.topic", "test")

		// Check execution order (high priority should execute first)
		expectedOrder := []int{2, 1, 3, 4}
		t.Assert(executionOrder, expectedOrder)
	})
}

func TestEvent_MultipleTopics(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()
		topic1Received := make(chan string, 1)
		topic2Received := make(chan string, 1)

		// Subscribe to multiple topics
		event.Subscribe("topic1", func(topic string, message any) {
			topic1Received <- message.(string)
		})

		event.Subscribe("topic2", func(topic string, message any) {
			topic2Received <- message.(string)
		})

		// Publish to both topics
		event.Publish("topic1", "Message for topic 1")
		event.Publish("topic2", "Message for topic 2")

		// Check received messages
		select {
		case msg := <-topic1Received:
			t.Assert(msg, "Message for topic 1")
		case <-time.After(time.Second):
			t.Error("Message for topic1 not received within timeout")
		}

		select {
		case msg := <-topic2Received:
			t.Assert(msg, "Message for topic 2")
		case <-time.After(time.Second):
			t.Error("Message for topic2 not received within timeout")
		}

		// Check topics list
		topics := event.Topics()
		t.Assert(len(topics), 2)
		t.AssertIN("topic1", topics)
		t.AssertIN("topic2", topics)
	})
}

func TestEvent_SubscribersCount(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		// Initially no subscribers
		t.Assert(event.SubscribersCount("count.topic"), 0)

		// Add subscribers
		event.Subscribe("count.topic", func(topic string, message any) {})
		t.Assert(event.SubscribersCount("count.topic"), 1)

		event.Subscribe("count.topic", func(topic string, message any) {}, gevent.PriorityHigh)
		t.Assert(event.SubscribersCount("count.topic"), 2)

		// Add subscriber to another topic
		event.Subscribe("another.topic", func(topic string, message any) {})
		t.Assert(event.SubscribersCount("count.topic"), 2)
		t.Assert(event.SubscribersCount("another.topic"), 1)

		// Unsubscribe one handler
		subscriber := event.Subscribe("count.topic", func(topic string, message any) {})
		t.Assert(event.SubscribersCount("count.topic"), 3)

		subscriber.Unsubscribe()
		t.Assert(event.SubscribersCount("count.topic"), 2)
	})
}

func TestEvent_Unsubscribe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()
		messageReceived := false

		// Subscribe and immediately unsubscribe
		subscriber := event.Subscribe("test.topic", func(topic string, message any) {
			messageReceived = true
		})

		subscriber.Unsubscribe()

		// Publish message
		event.Publish("test.topic", "test message")

		// Give some time for async processing
		time.Sleep(100 * time.Millisecond)

		// Message should not be received as subscriber was unsubscribed
		t.Assert(messageReceived, false)
	})
}

func TestEvent_UnsubscribeMultipleTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()
		messageReceived := false

		// Subscribe
		subscriber := event.Subscribe("test.topic", func(topic string, message any) {
			messageReceived = true
		})

		// Unsubscribe multiple times (should not panic)
		subscriber.Unsubscribe()
		subscriber.Unsubscribe()
		subscriber.Unsubscribe()

		// Publish message
		event.Publish("test.topic", "test message")

		// Give some time for async processing
		time.Sleep(100 * time.Millisecond)

		// Message should not be received
		t.Assert(messageReceived, false)
	})
}

func TestEvent_NoSubscribers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		event := gevent.New()

		// Publish to topic with no subscribers (should not panic)
		event.Publish("no.subscribers", "test message")
		event.PublishSync("no.subscribers", "test message")

		// Check subscribers count
		t.Assert(event.SubscribersCount("no.subscribers"), 0)
	})
}
