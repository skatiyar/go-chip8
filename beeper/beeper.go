package beeper

// typedef unsigned char Uint8;
// void AudioCallback(void *userdata, Uint8 *stream, int len);
import "C"

import (
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
