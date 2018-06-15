package cluster

import (
	"sort"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// PartitionConsumer allows code to consume individual partitions from the cluster.
//
// See docs for Consumer.Partitions() for more on how to implement this.
type PartitionConsumer interface {
	sarama.PartitionConsumer

	// Topic returns the consumed topic name
	Topic() string

	// Partition returns the consumed partition
	Partition() int32

	// InitialOffset returns the offset used for creating the PartitionConsumer instance.
	// The returned offset can be a literal offset, or OffsetNewest, or OffsetOldest
	InitialOffset() int64
  
	// MarkOffset marks the offset of a message as preocessed.
	MarkOffset(offset int64, metadata string)

	// ResetOffset resets the offset to a previously processed message.
	ResetOffset(offset int64, metadata string)
}

type partitionConsumer struct {
	sarama.PartitionConsumer

	state partitionState
	mu    sync.Mutex

	topic         string
	partition     int32
	initialOffset int64

	closeOnce sync.Once
	closeErr  error

	dying, dead chan none
}

func newPartitionConsumer(manager sarama.Consumer, topic string, partition int32, info offsetInfo, defaultOffset int64) (*partitionConsumer, error) {
	offset := info.NextOffset(defaultOffset)
	pcm, err := manager.ConsumePartition(topic, partition, offset)

	// Resume from default offset, if requested offset is out-of-range
	if err == sarama.ErrOffsetOutOfRange {
		info.Offset = -1
		offset = defaultOffset
		pcm, err = manager.ConsumePartition(topic, partition, offset)
	}
	if err != nil {
		return nil, err
	}

	return &partitionConsumer{
		PartitionConsumer: pcm,
		state:             partitionState{Info: info},

		topic:         topic,
		partition:     partition,
		initialOffset: offset,

		dying: make(chan none),
		dead:  make(chan none),
	}, nil
}

// Topic implements PartitionConsumer
func (c *partitionConsumer) Topic() string { return c.topic }

// Partition implements PartitionConsumer
func (c *partitionConsumer) Partition() int32 { return c.partition }

// InitialOffset implements PartitionConsumer
func (c *partitionConsumer) InitialOffset() int64 { return c.initialOffset }

// AsyncClose implements PartitionConsumer
func (c *partitionConsumer) AsyncClose() {
	c.closeOnce.Do(func() {
		c.closeErr = c.PartitionConsumer.Close()
		close(c.dying)
	})
}

// Close implements PartitionConsumer
func (c *partitionConsumer) Close() error {
	c.AsyncClose()
	<-c.dead
	return c.closeErr
}

func (c *partitionConsumer) waitFor(stopper <-chan none, errors chan<- error) {
	defer close(c.dead)

	for {
		select {
		case err, ok := <-c.Errors():
			if !ok {
				return
			}
			select {
			case errors <- err:
			case <-stopper:
				return
			case <-c.dying:
				return
			}
		case <-stopper:
			return
		case <-c.dying:
			return
		}
	}
}

func (c *partitionConsumer) multiplex(stopper <-chan none, messages chan<- *sarama.ConsumerMessage, errors chan<- error) {
	defer close(c.dead)

	for {
		select {
		case msg, ok := <-c.Messages():
			if !ok {
				return
			}
			select {
			case messages <- msg:
			case <-stopper:
				return
			case <-c.dying:
				return
			}
		case err, ok := <-c.Errors():
			if !ok {
				return
			}
			select {
			case errors <- err:
			case <-stopper:
				return
			case <-c.dying:
				return
			}
		case <-stopper:
			return
		case <-c.dying:
			return
		}
	}
}

func (c *partitionConsumer) getState() partitionState {
	c.mu.Lock()
	state := c.state
	c.mu.Unlock()

	return state
}

func (c *partitionConsumer) markCommitted(offset int64) {
	c.mu.Lock()
	if offset == c.state.Info.Offset {
		c.state.Dirty = false
	}
	c.mu.Unlock()
}

// MarkOffset implements PartitionConsumer
func (c *partitionConsumer) MarkOffset(offset int64, metadata string) {
	c.mu.Lock()
	if next := offset + 1; next > c.state.Info.Offset {
		c.state.Info.Offset = next
		c.state.Info.Metadata = metadata
		c.state.Dirty = true
	}
	c.mu.Unlock()
}

// ResetOffset implements PartitionConsumer
func (c *partitionConsumer) ResetOffset(offset int64, metadata string) {
	c.mu.Lock()
	if next := offset + 1; next <= c.state.Info.Offset {
		c.state.Info.Offset = next
		c.state.Info.Metadata = metadata
		c.state.Dirty = true
	}
	c.mu.Unlock()
}

// --------------------------------------------------------------------

type partitionState struct {
	Info       offsetInfo
	Dirty      bool
	LastCommit time.Time
}

// --------------------------------------------------------------------

type partitionMap struct {
	data map[topicPartition]*partitionConsumer
	mu   sync.RWMutex
}

func newPartitionMap() *partitionMap {
	return &partitionMap{
		data: make(map[topicPartition]*partitionConsumer),
	}
}

func (m *partitionMap) IsSubscribedTo(topic string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for tp := range m.data {
		if tp.Topic == topic {
			return true
		}
	}
	return false
}

func (m *partitionMap) Fetch(topic string, partition int32) *partitionConsumer {
	m.mu.RLock()
	pc, _ := m.data[topicPartition{topic, partition}]
	m.mu.RUnlock()
	return pc
}

func (m *partitionMap) Store(topic string, partition int32, pc *partitionConsumer) {
	m.mu.Lock()
	m.data[topicPartition{topic, partition}] = pc
	m.mu.Unlock()
}

func (m *partitionMap) Snapshot() map[topicPartition]partitionState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snap := make(map[topicPartition]partitionState, len(m.data))
	for tp, pc := range m.data {
		snap[tp] = pc.getState()
	}
	return snap
}

func (m *partitionMap) Stop() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var wg sync.WaitGroup
	for tp := range m.data {
		wg.Add(1)
		go func(p *partitionConsumer) {
			_ = p.Close()
			wg.Done()
		}(m.data[tp])
	}
	wg.Wait()
}

func (m *partitionMap) Clear() {
	m.mu.Lock()
	for tp := range m.data {
		delete(m.data, tp)
	}
	m.mu.Unlock()
}

func (m *partitionMap) Info() map[string][]int32 {
	info := make(map[string][]int32)
	m.mu.RLock()
	for tp := range m.data {
		info[tp.Topic] = append(info[tp.Topic], tp.Partition)
	}
	m.mu.RUnlock()

	for topic := range info {
		sort.Sort(int32Slice(info[topic]))
	}
	return info
}
