// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

// Package avutil is a utility library to aid portable multimedia programming.
// It contains safe portable string functions, random number generators, data structures,
// additional mathematics functions, cryptography and multimedia related functionality.
// Some generic features and utilities provided by the libavutil library
package avutil

/*
#cgo pkg-config: libavutil
#include <stdlib.h>
#include <libavutil/hwcontext.h>
#include <libavutil/hwcontext_qsv.h>
#include <libavcodec/avcodec.h>

#include <libavformat/avformat.h>
#include <libavutil/buffer.h>
#include <libavutil/error.h>

static int get_qsv_format(AVCodecContext* avctx, const enum AVPixelFormat* pix_fmts)
{
    while (*pix_fmts != AV_PIX_FMT_NONE) {
        if (*pix_fmts == AV_PIX_FMT_QSV) {
			AVBufferRef* hw_device_ref = avctx->opaque;
            AVHWFramesContext* frames_ctx;
            AVQSVFramesContext* frames_hwctx;
            int ret;

			avctx->hw_frames_ctx = av_hwframe_ctx_alloc(hw_device_ref);
			if (!avctx->hw_frames_ctx)
				return AV_PIX_FMT_NONE;
			frames_ctx = (AVHWFramesContext*)avctx->hw_frames_ctx->data;
			frames_hwctx = frames_ctx->hwctx;

			frames_ctx->format = AV_PIX_FMT_QSV;
			frames_ctx->sw_format = avctx->sw_pix_fmt;
			frames_ctx->width = FFALIGN(avctx->coded_width, 32);
			frames_ctx->height = FFALIGN(avctx->coded_height, 32);
			frames_ctx->initial_pool_size = 32;

			frames_hwctx->frame_type = MFX_MEMTYPE_VIDEO_MEMORY_DECODER_TARGET;
			ret = av_hwframe_ctx_init(avctx->hw_frames_ctx);
			if (ret < 0)
				return AV_PIX_FMT_NONE;

			return AV_PIX_FMT_QSV;
		}

		pix_fmts++;
	}

	fprintf(stderr, "The QSV pixel format not offered in get_format()\n");

	return AV_PIX_FMT_NONE;
}

void set_qvs_callback(AVCodecContext* ctx){
	ctx->get_format = get_qsv_format;
}
*/
import "C"
import (
	"unsafe"

	"github.com/giorgisio/goav/avcodec"
)

type (
	CodecContext     C.struct_AVCodecContext
	AVCodecHWConfig  C.struct_AVCodecHWConfig
	HWFramesContext  C.struct_AVHWFramesContext
	QSVFramesContext C.struct_AVQSVFramesContext
	AVHWDeviceType   C.enum_AVHWDeviceType
)

const (
	AV_HWDEVICE_TYPE_NONE                   = int(C.AV_HWDEVICE_TYPE_NONE)
	AV_HWDEVICE_TYPE_VDPAU                  = int(C.AV_HWDEVICE_TYPE_VDPAU)
	AV_HWDEVICE_TYPE_CUDA                   = int(C.AV_HWDEVICE_TYPE_CUDA)
	AV_HWDEVICE_TYPE_VAAPI                  = int(C.AV_HWDEVICE_TYPE_VAAPI)
	AV_HWDEVICE_TYPE_DXVA2                  = int(C.AV_HWDEVICE_TYPE_DXVA2)
	AV_HWDEVICE_TYPE_QSV                    = int(C.AV_HWDEVICE_TYPE_QSV)
	AV_HWDEVICE_TYPE_VIDEOTOOLBOX           = int(C.AV_HWDEVICE_TYPE_VIDEOTOOLBOX)
	AV_HWDEVICE_TYPE_D3D11VA                = int(C.AV_HWDEVICE_TYPE_D3D11VA)
	AV_HWDEVICE_TYPE_DRM                    = int(C.AV_HWDEVICE_TYPE_DRM)
	AV_HWDEVICE_TYPE_OPENCL                 = int(C.AV_HWDEVICE_TYPE_OPENCL)
	AV_HWDEVICE_TYPE_MEDIACODEC             = int(C.AV_HWDEVICE_TYPE_MEDIACODEC)
	AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX = int(C.AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX) //0x1
)

func AvHwdeviceFindTypeByName(name string) AVHWDeviceType {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	return (AVHWDeviceType)(C.av_hwdevice_find_type_by_name(Cname))
}

func AvCodecGetHwConfig(codec *avcodec.Codec, index int) *AVCodecHWConfig {
	config := (*AVCodecHWConfig)(C.avcodec_get_hw_config((*C.struct_AVCodec)(unsafe.Pointer(codec)), C.int(index)))
	return config
}

func (hw *AVCodecHWConfig) AvHwGetPixFormat() PixelFormat {
	return (PixelFormat)(hw.pix_fmt)
}

func (hw *AVCodecHWConfig) AvHwGetDeviceType() AVHWDeviceType {
	return (AVHWDeviceType)(hw.device_type)
}

func (hw *AVCodecHWConfig) AvHwGetMethods() int {
	return int(hw.methods)
}

func AvAllocBufferRef(buf *AvBufferRef) *AvBufferRef {
	return (*AvBufferRef)(C.av_buffer_ref((*C.struct_AVBufferRef)(unsafe.Pointer(buf))))
}

func AVHwDeviceCtxCreate(deviceCtx **AvBufferRef, dict *Dictionary, deviceType AVHWDeviceType, device string, flags int) int {
	if device == "" {
		return int(C.av_hwdevice_ctx_create((**C.struct_AVBufferRef)(unsafe.Pointer(deviceCtx)),
			(C.enum_AVHWDeviceType)(deviceType), nil,
			(*C.struct_AVDictionary)(unsafe.Pointer(dict)),
			C.int(flags)))
	}

	dev := C.CString(device)
	defer C.free(unsafe.Pointer(dev))

	ret := int(C.av_hwdevice_ctx_create((**C.struct_AVBufferRef)(unsafe.Pointer(deviceCtx)),
		(C.enum_AVHWDeviceType)(deviceType), dev,
		(*C.struct_AVDictionary)(unsafe.Pointer(dict)),
		C.int(flags)))

	return ret
}

// retrieve data from GPU to CPU
func AvHwFrameTransferData(dst, src *Frame, flag int) int {
	return int(C.av_hwframe_transfer_data((*C.struct_AVFrame)(unsafe.Pointer(dst)), (*C.struct_AVFrame)(unsafe.Pointer(src)), C.int(flag)))
}

func (b *AvBufferRef) HWFramesContext() *HWFramesContext {
	return (*HWFramesContext)((*C.struct_AVHWFramesContext)(unsafe.Pointer(b.data)))
}

func (c *CodecContext) SetQsvOpaque(ref *AvBufferRef) {
	c.opaque = unsafe.Pointer(ref)
}

func (c *CodecContext) SetQsvCallback() {
	C.set_qvs_callback((*C.struct_AVCodecContext)(unsafe.Pointer(c)))
}
