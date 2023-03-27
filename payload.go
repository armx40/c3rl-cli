package main

import (
	"fmt"
	pb "main/protobuf"
	"strings"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* custom error types */

/* */

type generalPayloadV2 struct {
	Payload interface{} `json:"payload" validate:"required"`
	Status  string      `json:"status" validate:"required"`
	Code    int         `json:"code" validate:"required"`
}

type userTokenPayload struct {
	Token     string             `json:"t"`
	Firstname string             `json:"f"`
	Lastname  string             `json:"l"`
	UserID    primitive.ObjectID `json:"ui"`
}

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

func (d *DeviceSettingsSurveyAnswerPayload) parse_sdcard_settings(settings *pb.SDCardSettings) (err error) {

	/* set operation mode */
	if ((settings.OPERATION_SETTINGS_SETTINGS >> SETTINGS_OPERATION_SETTINGS_MODE_BIT) & 0x1) == 1 {
		d.OperationMode = "Interrupt"
	} else {
		d.OperationMode = "WebSocket"
	}
	/**/

	/* set logging mode */
	if ((settings.OPERATION_SETTINGS_LOGGING_SETTINGS >> SETTINGS_LOGGING_SETTINGS_ENABLED) & 0x1) == 1 {
		d.LoggingEnable = true
	} else {
		d.LoggingEnable = false
	}
	/**/

	/* set sleep time */
	d.InterruptSleepTime = int(settings.OPERATION_SETTINGS_TIMINGS_SETTINGS>>SETTINGS_OPERATION_SETTINGS_TIMINGS_INTERRUPT_SLEEP_TIME_BIT) & 0xffff
	/**/

	/* set sleep time */
	d.InterruptOpenTime = int(settings.OPERATION_SETTINGS_TIMINGS_SETTINGS>>SETTINGS_OPERATION_SETTINGS_TIMINGS_INTERRUPT_OPEN_TIME_BIT) & 0xffff
	/**/

	/* set sensors */
	d.Sensors = nil

	if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_GNSS_ENABLED_BIT) & 0x1) == 1 {
		d.Sensors = append(d.Sensors, "GNSS")
	}

	if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_TEMPERATURE_ENABLED_BIT) & 0x1) == 1 {
		d.Sensors = append(d.Sensors, "Temperature")
	}

	if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_HUMIDITY_ENABLED_BIT) & 0x1) == 1 {
		d.Sensors = append(d.Sensors, "Humidity")
	}

	if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_TEMPERATURE_HUMIDITY_ENABLED_BIT) & 0x1) == 1 {
		d.Sensors = append(d.Sensors, "Temperature/Humidity")
	}

	if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_FLOWMETER_ENABLED_BIT) & 0x1) == 1 {
		d.Sensors = append(d.Sensors, "Flowmeter")
	}

	if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_MOTION_ENABLED_BIT) & 0x1) == 1 {
		d.Sensors = append(d.Sensors, "Motion")
	}

	// if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_MOTION_ENABLED_BIT)& 0x1) == 1 {
	// 	d.Sensors = append(d.Sensors, "Ambient Light")
	// }

	// if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_MOTION_ENABLED_BIT)& 0x1) == 1 {
	// 	d.Sensors = append(d.Sensors, "Inertial")
	// }

	// if ((settings.OPERATION_SETTINGS_SENSOR_DATA_SETTINGS >> SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_MOTION_ENABLED_BIT)& 0x1) == 1 {
	// 	d.Sensors = append(d.Sensors, "Magnetic")
	// }

	return
}

func (d *DeviceSettingsSurveyAnswerPayload) pretty_print() {

	key_color := color.New(color.FgHiWhite)
	value_color := color.New(color.FgHiGreen)
	red_color := color.New(color.FgHiRed)

	key_color.Printf("Operation Mode:")
	value_color.Printf(" %s", d.OperationMode)
	fmt.Printf("\n")

	key_color.Printf("Logging: ")
	if d.LoggingEnable {
		value_color.Printf("Enabled\n")
	} else {
		red_color.Printf("Disabled\n")
	}

	key_color.Printf("Interrupt Open Time:")
	value_color.Printf(" %d", d.InterruptOpenTime)
	fmt.Printf("s\n")

	key_color.Printf("Interrupt Sleep Time:")
	value_color.Printf(" %d", d.InterruptSleepTime)
	fmt.Printf("s\n")

	key_color.Printf("Sensors Enabled:")
	value_color.Printf(" %s", strings.Join(d.Sensors, ","))
	fmt.Printf("\n")

	return
}

type AuthLoginSurveyAnswerPayload struct {
	Username string
	Password string
}
