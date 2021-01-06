package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/kemalicecek/wav"
	"github.com/yobert/alsa"
)

var audioOutputDevice *alsa.Device

var isInitAudioOutput bool

var pause sync.Mutex

const (
	audioOutputPeriodSize = 4096
)

var (
	currentChannels      = 0
	currentRate          = 0
	currentBitsPerSample = -1
)

func main() {
	go func() {
		for {
			if err := initAudioOutput(); err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 3)
		}
	}()

	time.Sleep(time.Second * 1)

	wavFile, err := wav.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	wavFile.PrintWavHeader()

	go func() {
		err := audioOutputDeviceConfiguration(wavFile.GetNumChannels(), wavFile.GetSampleRate(), wavFile.GetBitsPerSample())
		if err != nil {
			fmt.Println(err.Error())
		}

		data := make([]byte, audioOutputPeriodSize*4)
		for i := 0; i < len(data); i++ {
			data[i] = 0
		}

		for {
			pause.Lock()
			if err := audioOutputDevice.Write(data, audioOutputPeriodSize); err != nil {
				return
			}
			pause.Unlock()
		}
	}()

	time.Sleep(time.Second * 2)

	err = playData(wavFile.Data)
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(time.Second * 1)

	fmt.Println("çal")

	// for {
	// 	beepDevice()
	// 	time.Sleep(time.Second * 2)
	// }

}

func initAudioOutput() error {
	if isInitAudioOutput {
		return nil
	}

	cards, err := alsa.OpenCards()
	if err != nil {
		return err
	}

	for _, card := range cards {
		devices, err := card.Devices()
		if err != nil {
			continue
		}
		for _, device := range devices {
			if device.Type != alsa.PCM || !device.Play {
				continue
			}

			audioOutputDevice = device

			isInitAudioOutput = true
			return nil
		}
	}

	return errors.New("audio device could not be installed")
}

func audioOutputDeviceConfiguration(channels, rate, bps int) error {

	if currentChannels != channels || currentRate != rate || currentBitsPerSample != bps {
		audioOutputDevice.Close()
	} else {
		return nil
	}

	err := audioOutputDevice.Open()
	if err != nil {
		return err
	}

	currentChannels = channels
	currentRate = rate
	currentBitsPerSample = bps

	_, err = audioOutputDevice.NegotiateChannels(currentChannels)
	if err != nil {
		return err
	}

	_, err = audioOutputDevice.NegotiateRate(currentRate)
	if err != nil {
		return err
	}

	bit := 0

	switch currentBitsPerSample {
	case 8:
		bit = 0
	case 16:
		bit = 2
	case 24:
		bit = 6
	case 32:
		bit = 10
	default:
		bit = 2

	}
	format, err := audioOutputDevice.NegotiateFormat(alsa.FormatType(bit))
	if err != nil {
		return err
	}
	fmt.Println(format)
	_, err = audioOutputDevice.NegotiatePeriodSize(audioOutputPeriodSize)
	if err != nil {
		return err
	}

	_, err = audioOutputDevice.NegotiateBufferSize(audioOutputPeriodSize * 2)
	if err != nil {
		return err
	}

	if err = audioOutputDevice.Prepare(); err != nil {
		return err
	}

	return nil
}

func playData(data []byte) error {
	if len(data) < (audioOutputPeriodSize * 4) {
		return errors.New("data is too short")
	}

	pause.Lock()
	defer func() { pause.Unlock() }()

	for i := 0; i < len(data)/(audioOutputPeriodSize*4)-1; i++ {
		if err := audioOutputDevice.Write(data[i*(audioOutputPeriodSize*4):(i+1)*(audioOutputPeriodSize*4)], audioOutputPeriodSize); err != nil {
			return err
		}
	}

	return nil
}

func beepDevice() error {

	pause.Lock()
	defer func() { pause.Unlock() }()

	duration := 350 * time.Millisecond
	for t := 0.; t < duration.Seconds(); {
		var buf bytes.Buffer

		for i := 0; i < 2048; i++ {
			v := math.Sin(t * 2 * math.Pi * 698) // A4
			v *= 0.3                             // make a little quieter

			sample := int16(v * math.MaxInt16)

			for c := 0; c < 2; c++ {
				binary.Write(&buf, binary.LittleEndian, sample)
			}

			t += 1 / float64(44100)
		}

		if err := audioOutputDevice.Write(buf.Bytes(), 2048); err != nil {
			return err
		}

	}
	// Wait for playback to complete.
	fmt.Printf("Playback should be complete now.\n")

	return nil
}
