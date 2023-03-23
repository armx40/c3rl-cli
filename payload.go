package main

import "fmt"

type LogLinePayload struct {
	LineLength  uint32
	LineOptions uint32
	LineTag     string
	LineCode    uint32
	LineTime    uint32
	LineLine    []byte
}

func (s *LogLinePayload) csv() (out []string) {
	out = append(out, fmt.Sprintf("%d", s.LineTime))
	out = append(out, s.LineTag)
	out = append(out, fmt.Sprintf("%d", s.LineCode))
	out = append(out, fmt.Sprintf("%s", string(s.LineLine)))

	return
}

func (s *LogLinePayload) csv_headers() (out []string) {
	out = append(out, "Time")
	out = append(out, "Tag")
	out = append(out, "Code")
	out = append(out, "Log")

	return
}

type DeviceSettingsPayload struct {
	OperationSettings           uint32
	OperationNumberSettings     uint32
	OperationTimingsSettings    uint32
	OperationSensorDataSettings uint32
	OperationLoggingSettings    uint32
}

type DeviceSettingsSurveyAnswerPayload struct {
	OperationMode      string
	Sensors            []string
	InterruptOpenTime  int
	InterruptSleepTime int
	LoggingEnable      bool
}

func (d *DeviceSettingsSurveyAnswerPayload) get_settings() (settings DeviceSettingsPayload) {

	/* set operation mode */
	operation_mode_websocket_bit := 1
	if d.OperationMode == "WebSocket" {
		operation_mode_websocket_bit = 0
	}
	settings.OperationSettings = uint32(helper_function_set_bit(uint64(settings.OperationSettings), SETTINGS_OPERATION_SETTINGS_MODE_BIT, uint8(operation_mode_websocket_bit)))
	/**/

	/* set logging mode */
	operation_mode_logging_bit := 0
	if d.LoggingEnable {
		operation_mode_logging_bit = 1
	}
	settings.OperationLoggingSettings = uint32(helper_function_set_bit(uint64(settings.OperationLoggingSettings), SETTINGS_LOGGING_SETTINGS_ENABLED, uint8(operation_mode_logging_bit)))

	/**/

	/* set sleep time */
	settings.OperationTimingsSettings = uint32(helper_function_set_bits(uint64(settings.OperationTimingsSettings), SETTINGS_OPERATION_SETTINGS_TIMINGS_INTERRUPT_SLEEP_TIME_BIT, uint64(d.InterruptSleepTime), 16))
	/**/

	/* set open time */
	settings.OperationTimingsSettings = uint32(helper_function_set_bits(uint64(settings.OperationTimingsSettings), SETTINGS_OPERATION_SETTINGS_TIMINGS_INTERRUPT_OPEN_TIME_BIT, uint64(d.InterruptOpenTime), 16))
	/**/

	/* set sensors */
	for i := range d.Sensors {
		if d.Sensors[i] == "GNSS" {
			settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_GNSS_ENABLED_BIT, 1))
		}

		if d.Sensors[i] == "Temperature" {
			settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_TEMPERATURE_ENABLED_BIT, 1))
		}
		if d.Sensors[i] == "Humidity" {
			settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_HUMIDITY_ENABLED_BIT, 1))
		}
		if d.Sensors[i] == "Temperature/Humidity" {
			settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_TEMPERATURE_HUMIDITY_ENABLED_BIT, 1))
		}

		if d.Sensors[i] == "Flowmeter" {
			settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_FLOWMETER_ENABLED_BIT, 1))
		}

		if d.Sensors[i] == "Motion" {
			settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_MOTION_ENABLED_BIT, 1))
		}

		if d.Sensors[i] == "Ambient Light" {
			// settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_Ambi, 1))
		}
		if d.Sensors[i] == "Inertial" {
			// settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_GNSS_ENABLED_BIT, 1))
		}
		if d.Sensors[i] == "Magnetic" {
			// settings.OperationSensorDataSettings = uint32(helper_function_set_bit(uint64(settings.OperationSensorDataSettings), SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_GNSS_ENABLED_BIT, 1))
		}
	}
	/**/
	return
}
