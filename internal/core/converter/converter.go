package converter

import (
	"context"
	"fmt"
	"simple-media-server/pkg/ffmpeg"
	"sync"

	"go.uber.org/zap"
)

type item struct {
	From string
	To   string
}

type Convert struct {
	log    *zap.Logger
	ffmpeg ffmpeg.FFMpeg
	state  map[string]bool
	queueC chan item
	mu     *sync.Mutex
}

func New(ctx context.Context, log *zap.Logger) (*Convert, error) {
	c := &Convert{
		log:    log,
		state:  map[string]bool{},
		queueC: make(chan item, 10),
		mu:     &sync.Mutex{},
	}

	ffmpg, err := ffmpeg.New()
	if err != nil {
		c.log.Error("failed init ffmpeg", zap.Error(err))
		return nil, err
	}

	c.ffmpeg = *ffmpg

	go c.worker(ctx)

	return c, nil
}

func (c *Convert) AviToMP4Q(from, to string) error {
	select {
	case c.queueC <- item{From: from, To: to}:
		c.state[from] = false

	default:
		c.log.Info("queue is full")
		return fmt.Errorf("queue is full")

	}

	return nil
}

func (c *Convert) Status(from string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ok := c.state[from]; ok {
		return true
	}

	return false
}

func (c *Convert) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.log.Info("converter worker stopped")
			return
		case d := <-c.queueC:
			c.mu.Lock()
			c.state[d.From] = true
			c.mu.Unlock()

			c.log.Debug("start converting", zap.String("from", d.From), zap.String("to", d.To))
			output, err := c.ffmpeg.AviToMP4(d.From, d.To)
			if err != nil {
				c.log.Error("failed convert", zap.Error(err))
			}
			c.mu.Lock()
			delete(c.state, d.From)
			c.mu.Unlock()

			c.log.Debug("convert avi to mp4", zap.String("from", d.From), zap.String("to", d.To), zap.String("output", output))
		}
	}
}
