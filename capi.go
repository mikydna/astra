package astra

// fix: C-FLAGS with  ENV vars?

/*
#cgo CFLAGS: -I/Users/andy/Desktop/AstraSDK-0.5.0-20160426T102621Z-darwin-x64/include
#cgo LDFLAGS: -Wl,-rpath,/Users/andy/Desktop/AstraSDK-0.5.0-20160426T102621Z-darwin-x64/lib -lastra -lastra_core -lastra_core_api
#include <stdlib.h>
#include <astra/capi/astra.h>
*/
import "C"

import (
	"errors"
	"log"
)

type Status uint8

const (
	StatusSuccess Status = iota
	StatusInvalidParameter
	StatusDeviceError
	StatusTimeout
	StatusInvalidParameterToken
	StatusInvalidOperation
	StatusInternalError
	StatusUninitialized
)

func (s Status) String() string {
	var str string
	switch s {
	case StatusSuccess:
		str = "SUCCESS"
	case StatusInvalidParameter:
		str = "INVALID_PARAMETER"
	case StatusTimeout:
		str = "TIMEOUT"
	case StatusInvalidParameterToken:
		str = "INVALID_PARAMATER_TOKEN"
	case StatusInvalidOperation:
		str = "INVALID_OPERATION"
	case StatusInternalError:
		str = "INTERNAL_ERROR"
	case StatusUninitialized:
		str = "UNINITIALIZED"
	}

	return str
}

func (s Status) Error() error {
	var err error
	switch s {
	case StatusSuccess:
		err = nil
	default:
		err = errors.New(s.String())
	}

	return err
}

type StreamSetConnection C.astra_streamsetconnection_t

type Reader C.astra_reader_t
type DepthStream C.astra_depthstream_t
type HandStream C.astra_handstream_t

type ReaderFrame C.astra_reader_frame_t
type FrameIndex C.astra_frame_index_t

type DepthFrame C.astra_depthframe_t
type HandFrame C.astra_handframe_t

// ??
// type ImageMetadata C.astra_image_metadata_t

func Initialize() Status {
	defer func() {
		if err := recover(); err != nil {
			log.Println("fix required", err)
		}
	}()

	rc := C.astra_initialize()

	return Status(rc)
}

func Terminate() Status {
	rc := C.astra_terminate()

	return Status(rc)
}

func OpenStream(deviceAddr string, conn *StreamSetConnection) Status {
	addr := C.CString(deviceAddr)
	rc := C.astra_streamset_open(addr, conn)
	return Status(rc)
}

func CloseStream(conn *StreamSetConnection) Status {
	rc := C.astra_streamset_close(conn)
	return Status(rc)
}

func CreateReader(conn StreamSetConnection, reader *Reader) Status {
	rc := C.astra_reader_create(conn, reader)
	return Status(rc)
}

func DestroyReader(reader *Reader) Status {
	rc := C.astra_reader_destroy(reader)
	return Status(rc)
}

func Update() Status {
	rc := C.astra_temp_update()
	return Status(rc)
}

func OpenReaderFrame(reader Reader, frame *ReaderFrame) Status {
	rc := C.astra_reader_open_frame(reader, 0, frame)
	return Status(rc)
}

func CloseReaderFrame(frame *ReaderFrame) Status {
	rc := C.astra_reader_close_frame(frame)
	return Status(rc)
}

/* depth stream */

func StartDepthStream(depthStream DepthStream) Status {
	rc := C.astra_stream_start(depthStream)
	return Status(rc)
}

func GetDepthStream(reader Reader, depthStream *DepthStream) Status {
	rc := C.astra_reader_get_depthstream(reader, depthStream)
	return Status(rc)
}

func GetDepthStreamFOV(depthStream DepthStream) (float32, float32, Status) {
	var hfov, vfov C.float

	if rc := C.astra_depthstream_get_hfov(depthStream, &hfov); Status(rc) != StatusSuccess {
		return -1, -1, Status(rc)
	}

	if rc := C.astra_depthstream_get_vfov(depthStream, &vfov); Status(rc) != StatusSuccess {
		return -1, -1, Status(rc)
	}

	return float32(hfov), float32(vfov), StatusSuccess
}

func GetDepthFrame(frame ReaderFrame, depthFrame *DepthFrame) (int, Status) {
	if rc := C.astra_frame_get_depthframe(frame, depthFrame); Status(rc) != StatusSuccess {
		return -1, Status(rc)
	}

	depthFrameIndex := new(C.astra_frame_index_t)
	if rc := C.astra_depthframe_get_frameindex(*depthFrame, depthFrameIndex); Status(rc) != StatusSuccess {
		return -1, Status(rc)
	}

	return int((C.int)(*depthFrameIndex)), StatusSuccess
}

func GetDepthFrameBuffer(depthFrame DepthFrame) ([]int16, Status) {
	blen := new(C.size_t)
	if rc := C.astra_depthframe_get_data_byte_length(depthFrame, blen); Status(rc) != StatusSuccess {
		return nil, Status(rc)
	}

	buffer := make([]int16, (C.int)(*blen)/2)
	if rc := C.astra_depthframe_copy_data(depthFrame, (*C.int16_t)(&buffer[0])); Status(rc) != StatusSuccess {
		return nil, Status(rc)
	}

	return buffer, StatusSuccess
}

func GetDepthFrameMetadata(depthFrame DepthFrame) (int, int, Status) {
	metadata := new(C.astra_image_metadata_t)
	if rc := C.astra_depthframe_get_metadata(depthFrame, metadata); Status(rc) != StatusSuccess {
		return 0, 0, Status(rc)
	}

	return int((C.uint)(metadata.width)), int((C.uint)(metadata.height)), StatusSuccess
}
