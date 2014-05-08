package pulse

/*
#include "dde-pulse.h"
*/
import "C"

type SinkInput struct {
	Index       uint32
	Name        string
	OwnerModule uint32
	Client      uint32
	Sink        uint32

	//sample_spec

	ChannelMap ChannelMap
	Volume     CVolume

	//buffer usec
	//sink usec

	ResampleMethod string
	Driver         string

	Mute     bool
	PropList map[string]string
	Corked   int

	HasVolume      bool
	VolumeWritable bool

	//format
}

func (sink *SinkInput) GetAvgVolume() float64 {
	return sink.Volume.Avg()
}

func (sink *SinkInput) SetAvgVolume(v float64) {
	sink.Volume.SetAvg(v, sink.ChannelMap)
	C.pa_context_set_sink_input_volume(GetContext().ctx, C.uint32_t(sink.Index), &sink.Volume.core, C.success_cb, nil)
}

func (sink *SinkInput) SetMute(mute bool) {
	_mute := 0
	if mute {
		_mute = 1
	}
	C.pa_context_set_sink_input_mute(GetContext().ctx, C.uint32_t(sink.Index), C.int(_mute), C.success_cb, nil)
}

func (sink *SinkInput) SetBalance(balance float64) {
	sink.Volume.SetBalance(sink.ChannelMap, balance)
	C.pa_context_set_sink_input_volume(GetContext().ctx, C.uint32_t(sink.Index), &sink.Volume.core, C.success_cb, nil)
}
