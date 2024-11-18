package common

import (
	"NanoKVM-Server/config"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	kvmVision     *KvmVision
	kvmVisionOnce sync.Once
)

type KvmVision struct {
	mutex sync.Mutex
}

func GetKvmVision() *KvmVision {
	kvmVisionOnce.Do(func() {
		kvmVision = &KvmVision{}

		conf := config.GetInstance()
		logLevel := strings.ToLower(conf.Logger.Level)

		logEnable := false
		if logLevel == "debug" {
			logEnable = true
		}

		log.Debugf("TODO(vially): initialize kvm vision %t", logEnable)

		log.Debugf("kvm vision initialized")
	})

	return kvmVision
}

func (k *KvmVision) ReadMjpeg(width uint16, height uint16, quality uint16) (data []byte, result int) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	log.Debugf("TODO(vially): implement `KvmVision.ReadMjpeg`")

	result = -1
	if result < 0 {
		log.Errorf("failed to read kvm image: %v", result)
		return
	}

	log.Debugf("read kvm image: %v", result)
	return
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

	log.Debugf("TODO(vially): implement `KvmVision.Close`")

	log.Debugf("stop kvm vision...")
}
