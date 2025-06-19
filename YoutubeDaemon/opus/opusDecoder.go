package opus

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed opus.wasm
var bin []byte

var (
	Mod     api.Module
	Runtime wazero.Runtime
	malloc, free,
	opus_strerror,
	opus_decode,
	opus_decode_float,
	decoder_create api.Function
)

func init() {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	Runtime = r

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	var err error
	Mod, err = r.InstantiateWithConfig(ctx, bin,
		wazero.NewModuleConfig().WithStartFunctions("_initialize"))
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	decoder_create = Mod.ExportedFunction("opus_decoder_create")
	malloc = Mod.ExportedFunction("malloc")
	free = Mod.ExportedFunction("free")
	opus_strerror = Mod.ExportedFunction("opus_strerror")
	opus_decode = Mod.ExportedFunction("opus_decode")
	opus_decode_float = Mod.ExportedFunction("opus_decode_float")
}

func main() {
	ptr, err := NewDecoder(48000, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoder pointer: 0x%X\n", ptr)
}

type Decoder struct {
	ptr      uintptr // pointer to OpusDecoder in WASM memory
	channels int     // number of output channels (e.g. 1=mono, 2=stereo)
}

func NewDecoder(sample_rate int, channels int) (Decoder, error) {
	ctx := context.Background()

	offset := Malloc(4)
	defer Free(offset)
	// create decoder, which returns pointer to decoder
	result, err := decoder_create.Call(ctx, uint64(sample_rate), uint64(channels), uint64(offset))
	if err != nil {
		return Decoder{}, fmt.Errorf("WASM call failed: %w", err)
	}

	errorCode, ok := Mod.Memory().ReadUint32Le(uint32(offset))
	if !ok {
		return Decoder{}, errors.New("could not read error code from WASM memory")
	}
	if errorCode != 0 {
		return Decoder{}, fmt.Errorf("opus error in NewDecoder: %s", OpusStrerror(int32(errorCode)))
	}

	return Decoder{ptr: uintptr((result[0])), channels: channels}, nil
}

func Malloc(bytes int) uintptr {
	res, err := malloc.Call(context.Background(), uint64(bytes))
	if err != nil {
		panic(fmt.Errorf("malloc failed: %w", err))
	}
	return uintptr(res[0])
}

func Free(offset uintptr) {
	_, err := free.Call(context.Background(), uint64(offset))
	if err != nil {
		panic(fmt.Errorf("free failed: %w", err))
	}
}
func OpusStrerror(code int32) string {
	ctx := context.Background()

	res, err := opus_strerror.Call(ctx, uint64(uint32(code)))
	if err != nil {
		return fmt.Sprintf("Opus error %d (failed to get string): %v", code, err)
	}

	ptr := uint32(res[0])
	mem := Mod.Memory()

	// Read null-terminated string from WASM memory
	bytes := []byte{}
	for {
		b, ok := mem.ReadByte(ptr)
		if !ok {
			break
		}
		if b == 0 {
			break
		}
		bytes = append(bytes, b)
		ptr++
	}
	return string(bytes)
}

func (d *Decoder) DecodeFloat32(data []byte, pcm []float32) (int, error) {
	if d.ptr == 0 {
		return 0, errors.New("opus: decoder uninitialized")
	}
	if len(data) == 0 {
		return 0, errors.New("opus: no data supplied")
	}
	if len(pcm) == 0 {
		return 0, errors.New("opus: target buffer empty")
	}
	if cap(pcm)%d.channels != 0 {
		return 0, errors.New("opus: target buffer capacity must be multiple of channels")
	}

	ctx := context.Background()
	mem := Mod.Memory()

	// Allocate and write encoded data to WASM memory
	dataPtr := Malloc(len(data))
	defer Free(dataPtr)
	if !mem.Write(uint32(dataPtr), data) {
		return 0, errors.New("failed to write input data")
	}

	// Allocate output buffer in WASM memory
	frameSize := cap(pcm) / d.channels
	pcmBytes := len(pcm) * 4 // float32 = 4 bytes
	pcmPtr := Malloc(pcmBytes)
	defer Free(pcmPtr)

	// Call opus_decode_float
	result, err := opus_decode_float.Call(ctx,
		uint64(d.ptr),
		uint64(dataPtr),
		uint64(len(data)),
		uint64(pcmPtr),
		uint64(frameSize),
		uint64(0), // no FEC
	)
	if err != nil {
		return 0, fmt.Errorf("opus_decode_float call failed: %w", err)
	}

	samples := int(result[0])

	// Read float32 samples from WASM memory
	raw, ok := mem.Read(uint32(pcmPtr), uint32(samples*d.channels*4))
	if !ok || len(raw) < samples*d.channels*4 {
		return 0, errors.New("failed to read float32 PCM output")
	}

	// Decode little-endian float32s into Go slice
	for i := 0; i < len(raw); i += 4 {
		bits := uint32(raw[i]) | uint32(raw[i+1])<<8 | uint32(raw[i+2])<<16 | uint32(raw[i+3])<<24
		pcm[i/4] = math.Float32frombits(bits)
	}

	return int(samples), nil
}
