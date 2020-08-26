package wav

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path"
)

//Header struct for initilize headers
type Header struct {
	Path          *os.File
	Title         string
	Duration      float32
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

//Init wav file
func Init(filePath string) (*Header, error) {
	// _, err := os.Stat(filePath)
	// if os.IsNotExist(err) {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	wavFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 45)
	_, err = wavFile.Read(b)
	for i := 0; i < 45; i++ {
		fmt.Printf("%d=%x ", i, b[i])

	}
	fmt.Println()
	var wavHeader Header
	wavHeader.Path = wavFile
	wavHeader.Title = path.Base(filePath)
	wavHeader.chunkID = b[0:4]         //string
	wavHeader.chunkSize = b[4:8]       //LE Uint32
	wavHeader.format = b[8:12]         //string
	wavHeader.subchunk1ID = b[12:16]   //string
	wavHeader.subchunk1Size = b[16:20] //LE Uint32
	wavHeader.audioFormat = b[20:22]   //LE Uint16
	wavHeader.numChannels = b[22:24]   //LE Uint16
	wavHeader.sampleRate = b[24:28]    //LE Uint32
	wavHeader.byteRate = b[28:32]      //LE Uint32
	wavHeader.blockAlign = b[32:34]    //LE Uint16
	wavHeader.bitsPerSample = b[34:36] //LE Uint16
	wavHeader.Subchunk2ID = b[36:40]   //string
	wavHeader.Subchunk2Size = b[40:45] //LE Uint32
	wavHeader.Duration = float32(binary.LittleEndian.Uint32(wavHeader.Subchunk2Size)) / float32(binary.LittleEndian.Uint32(wavHeader.byteRate))

	return &wavHeader, err
}

//PrintWavHeader ...
func (wavHeader *Header) PrintWavHeader() error {
	fmt.Println("\t\t\tWAVE FILE HEADER\t\t\t")
	fmt.Println("==============================================================")
	fmt.Printf("Title = %s, ", wavHeader.Title)
	fmt.Printf("Duration = %.2f, ", wavHeader.Duration)
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
func (wavHeader *Header) GetChunkSize() uint32 {
	return binary.LittleEndian.Uint32(wavHeader.chunkSize)
}

//GetFormat ...
func (wavHeader *Header) GetFormat() string {
	return string(wavHeader.format)
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
