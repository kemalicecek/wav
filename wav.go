package wav

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

type File struct {
	Data []byte
	Header
}

//Header struct for initilize headers
type Header struct {
	Title         string
	Duration      uint32
	chunkID       []byte
	chunkSize     []byte
	format        []byte
	subchunk1ID   []byte
	subchunk1Size []byte
	audioFormat   []byte
	numChannels   []byte
	sampleRate    []byte
	byteRate      []byte
	blockAlign    []byte
	bitsPerSample []byte
	Subchunk2ID   []byte
	Subchunk2Size []byte
}

//Open wav file
func Open(filePath string) (*File, error) {
	// _, err := os.Stat(filePath)
	// if os.IsNotExist(err) {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	wavFileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var wavFile File

	if len(wavFileBytes) < 44 {
		return nil, errors.New("corrupted file")
	}
	b := wavFileBytes[:80]

	wavFile.Data = wavFileBytes[80:]

	wavFile.Title = path.Base(filePath)
	wavFile.chunkID = b[0:4]         //string
	wavFile.chunkSize = b[4:8]       //LE Uint32
	wavFile.format = b[8:12]         //string
	wavFile.subchunk1ID = b[12:16]   //string
	wavFile.subchunk1Size = b[16:20] //LE Uint32

	switch string(wavFile.subchunk1ID) {
	case "fmt ":
		wavFile.audioFormat = b[20:22]   //LE Uint16
		wavFile.numChannels = b[22:24]   //LE Uint16
		wavFile.sampleRate = b[24:28]    //LE Uint32
		wavFile.byteRate = b[28:32]      //LE Uint32
		wavFile.blockAlign = b[32:34]    //LE Uint16
		wavFile.bitsPerSample = b[34:36] //LE Uint16
		wavFile.Subchunk2ID = b[36:40]   //string
		switch string(b[36:40]) {
		case "data":
			wavFile.Subchunk2Size = b[40:44] //LE Uint32
		case "LIST":
			wavFile.Subchunk2Size = b[74:78] //LE Uint32
			// default:
			// 	wavFile.Subchunk2Size = [0] //LE Uint32
		}

	case "JUNK":
		wavFile.audioFormat = b[56:58]   //LE Uint16
		wavFile.numChannels = b[38:40]   //LE Uint16
		wavFile.sampleRate = b[52:56]    //LE Uint32
		wavFile.byteRate = b[64:68]      //LE Uint32
		wavFile.blockAlign = b[68:70]    //LE Uint16
		wavFile.bitsPerSample = b[60:63] //LE Uint16
		wavFile.Subchunk2ID = b[72:76]   //string
		wavFile.Subchunk2Size = b[76:80] //LE Uint32
	default:
		return nil, errors.New(string(wavFileBytes[12:16]) + " is an unsupported file")
	}
	wavFile.Duration = binary.LittleEndian.Uint32(wavFile.Subchunk2Size) / binary.LittleEndian.Uint32(wavFile.byteRate)

	return &wavFile, nil
}

//PrintWavHeader ...
func (wavHeader *Header) PrintWavHeader() error {
	fmt.Println("\t\t\tWAVE FILE HEADER\t\t\t")
	fmt.Println("==============================================================")
	fmt.Printf("Title = %s, ", wavHeader.Title)
	fmt.Printf("Duration = %d:%d, ", (wavHeader.Duration / 60), (wavHeader.Duration % 60))
	fmt.Printf("ChunkID = %s, ", wavHeader.chunkID)
	fmt.Printf("Chunk Size = %d, ", binary.LittleEndian.Uint32(wavHeader.chunkSize))
	fmt.Printf("Format = %s\n", wavHeader.format)
	fmt.Printf("Subchunk1ID = %s, ", wavHeader.subchunk1ID)
	fmt.Printf("Subchunk1Size = %d, ", binary.LittleEndian.Uint32(wavHeader.subchunk1Size))
	fmt.Printf("AudioFormat = %d, ", binary.LittleEndian.Uint16(wavHeader.audioFormat))
	fmt.Printf("NumChannels = %d,\n", binary.LittleEndian.Uint16(wavHeader.numChannels))
	fmt.Printf("SampleRate = %d, ", binary.LittleEndian.Uint32(wavHeader.sampleRate))
	fmt.Printf("ByteRate = %d, ", binary.LittleEndian.Uint32(wavHeader.byteRate))
	fmt.Printf("BlockAlign = %d, ", binary.LittleEndian.Uint16(wavHeader.blockAlign))
	fmt.Printf("BitsPerSample = %d\n", binary.LittleEndian.Uint16(wavHeader.bitsPerSample))
	// extraParamSize := b[40:44]
	// fmt.Printf("BlockAlign = %s\n", extraParamSize)
	// extraParams := b[20:22]
	// fmt.Printf("ExtraParams = %s\n", extraParams)
	fmt.Printf("Subchunk2ID = %s, ", wavHeader.Subchunk2ID)
	fmt.Printf("Subchunk2Size = %d\n", binary.LittleEndian.Uint32(wavHeader.Subchunk2Size))
	return nil
}

//GetChunkID ...
func (wavHeader *Header) GetChunkID() string {
	return string(wavHeader.chunkID)
}

//GetChunkSize ...
func (wavHeader *Header) GetChunkSize() int {
	return int(binary.LittleEndian.Uint32(wavHeader.chunkSize))
}

//GetFormat ...
func (wavHeader *Header) GetFormat() string {
	return string(wavHeader.format)
}

//GetAudioFormat ...
func (wavHeader *Header) GetAudioFormat() int {
	return int(binary.LittleEndian.Uint16(wavHeader.audioFormat))
}

//GetNumChannels ...
func (wavHeader *Header) GetNumChannels() int {
	return int(binary.LittleEndian.Uint16(wavHeader.numChannels))
}

//GetSampleRate ...
func (wavHeader *Header) GetSampleRate() int {
	return int(binary.LittleEndian.Uint32(wavHeader.sampleRate))
}

//GetBitsPerSample ...
func (wavHeader *Header) GetBitsPerSample() int {
	return int(binary.LittleEndian.Uint16(wavHeader.bitsPerSample))
}

//IsWav ...
func (wavHeader *Header) IsWav() error {
	if string(wavHeader.format) == "WAVE" {
		return nil
	}
	err := errors.New("This is not a WAV file")
	return err
}

//Close function
func (wavHeader *Header) Close() {

}
