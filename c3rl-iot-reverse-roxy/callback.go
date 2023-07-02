package c3rliotroxy

import "time"

var callback_timeout_timer *time.Ticker

const callback_timeout_seconds = 5

func callback_init() {
	callback_timeout_timer = time.NewTicker(callback_timeout_seconds * time.Second)

	go func() {
		for {
			select {
			case _ = <-callback_timeout_timer.C:
				callback_send_message(ROXY_CALLBACK_TIMEOUT)
				return
			}
		}
	}()

}

func Callback_deinit() {
	if callback_timeout_timer != nil {
		callback_timeout_timer.Stop()
	}

}

func callback_send_message(callback_message int) {

	callback_reset_timeout()

	if main_app_callback_channel != nil {
		*main_app_callback_channel <- callback_message
	}

}

func callback_reset_timeout() {
	if callback_timeout_timer != nil {
		callback_timeout_timer.Reset(callback_timeout_seconds * time.Second)
	}
}
