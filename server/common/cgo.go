package common

import (
	"NanoKVM-Server/config"
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	kvmVision     *KvmVision
	kvmVisionOnce sync.Once
)

type KvmVision struct {
	mutex                 sync.Mutex
	dev                   *device.Device
	cancel                context.CancelFunc
	width                 uint16
	height                uint16
	skipStreamStartFrames int
	currentFrameIndex     int
}

func GetKvmVision() *KvmVision {
	kvmVisionOnce.Do(func() {
		conf := config.GetInstance()
		kvmVision = &KvmVision{
			width:                 0,
			height:                0,
			currentFrameIndex:     0,
			skipStreamStartFrames: conf.SkipStreamStartFrames,
		}

		log.Debugf("kvm vision initialized")
	})

	return kvmVision
}

func (k *KvmVision) ReadMjpeg(width uint16, height uint16, quality uint16) (data []byte, result int) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if width != k.width || height != k.height || k.dev == nil {
		log.Infof("updating stream pixel format: %dx%d (from %dx%d), quality: %d", width, height, k.width, k.height, quality)
		if err := k.updateMjpegPixelFormat(width, height); err != nil {
			log.Errorf("failed to update pixel format dimentions (%dx%d): %v", width, height, err)
			result = -1
			return
		}
	}

	data = <-k.dev.GetOutput()
	if k.currentFrameIndex < k.skipStreamStartFrames {
		k.currentFrameIndex++
		log.Warnf("skipping frame %d", k.currentFrameIndex)
		result = -1
		return
	}

	if len(data) == 0 {
		log.Warnf("v4l2 device returned empty data")
		result = -1
		return
	}

	log.Debugf("read kvm image: %v", result)
	return
}

func (k *KvmVision) updateMjpegPixelFormat(width uint16, height uint16) error {
	k.closeStream()

	options := []device.Option{device.WithBufferSize(1)}
	if width > 0 && height > 0 {
		options = append(options, device.WithPixFormat(v4l2.PixFormat{
			Width:       uint32(width),
			Height:      uint32(height),
			PixelFormat: v4l2.PixelFmtMJPEG,
			Field:       v4l2.FieldNone,
		}))
	}

	conf := config.GetInstance()
	dev, err := device.Open(conf.VideoDevice, options...)
	if err != nil {
		log.Errorf("v4l2 open failed: %v", err)
		return err
	}

	ctx, cancel := context.WithCancel(context.TODO())
	k.cancel = cancel
	if err := dev.Start(ctx); err != nil {
		log.Errorf("v4l2 start failed: %v", err)
		return err
	}

	k.dev = dev
	k.width = width
	k.height = height
	k.currentFrameIndex = 0
	return nil
}

func (k *KvmVision) closeStream() {
	if k.dev == nil {
		return
	}

	k.cancel()
	if err := k.dev.Close(); err != nil {
		log.Errorf("close failed: %v", err)
	}

	k.dev = nil
	k.width = 0
	k.height = 0
	k.currentFrameIndex = 0
}

func (k *KvmVision) ReadH264(width uint16, height uint16, bitRate uint16) (data []byte, sps []byte, pps []byte, result int) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	log.Debugf("TODO(vially): implement `KvmVision.ReadH264`")

	result = -1
	if result < 0 {
		log.Errorf("failed to read kvm image: %v", result)
		return
	}

	log.Debugf("read kvm image: %v", result)
	return
}

func (k *KvmVision) Close() {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.closeStream()
	log.Debugf("stop kvm vision...")
}
