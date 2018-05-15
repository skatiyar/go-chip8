package beeper

// typedef unsigned char Uint8;
// void AudioCallback(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"math"
	"reflect"
	"time"
	"unsafe"

	sdl "github.com/veandco/go-sdl2/sdl"
)

type Beeper struct {
	deviceId sdl.AudioDeviceID
}

func Init() (Beeper, error) {
	instance := Beeper{}

	desiredSpec := sdl.AudioSpec{
		Freq:     44100,
		Format:   sdl.AUDIO_S16SYS,
		Channels: 1,
		Samples:  2048,
		Callback: sdl.AudioCallback(C.AudioCallback),
	}
	obtainedSpec := sdl.AudioSpec{}

	deviceId, deviceErr := sdl.OpenAudioDevice(sdl.GetAudioDeviceName(0, false), false, &desiredSpec, &obtainedSpec, 0)
	if deviceErr != nil {
		return instance, deviceErr
	}

	instance.deviceId = deviceId

	return instance, nil
}

//export AudioCallback
func AudioCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	buf := *(*[]C.Uint8)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(stream)),
		Len:  n,
		Cap:  n,
	}))

	var phase float64
	for i := 0; i < n; i += 2 {
		phase += 2 * math.Pi * 440 / 44100
		sample := C.Uint8((math.Sin(phase) + 0.999999) * 128)
		buf[i] = sample
		buf[i+1] = sample
	}
}

func (b *Beeper) Play() {
	sdl.PauseAudioDevice(b.deviceId, false)
	go func() {
		dTicker := time.NewTimer(time.Second / 10)
		select {
		case <-dTicker.C:
			sdl.PauseAudioDevice(b.deviceId, true)
		}
	}()
}

func (b *Beeper) Close() {
	sdl.CloseAudioDevice(b.deviceId)
}
