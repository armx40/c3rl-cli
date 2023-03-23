package main

import (
	"fmt"
	"log"
	"os"

	pb "main/protobuf"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"

	"github.com/golang/protobuf/proto"
)

const (
	SETTINGS_OPERATION_SETTINGS_MODE_BIT                             = 0
	SETTINGS_OPERATION_SETTINGS_DONT_SLEEP_ON_INTERRUPT_BIT          = 1
	SETTINGS_OPERATION_SETTINGS_FETCH_SETTINGS_ON_POWER_ON_BIT       = 2
	SETTINGS_OPERATION_SETTINGS_RETRY_INTERRUPT_BIT                  = 3
	SETTINGS_OPERATION_SETTINGS_DONT_FETCH_SETTINGS_ON_INTERRUPT_BIT = 4

	SETTINGS_OPERATION_SETTINGS_NUMBERS_WEBSOCKET_RETRIES_BIT   = 0
	SETTINGS_OPERATION_SETTINGS_NUMBERS_INTERRUPTS_RETRIES_BIT  = 8
	SETTINGS_OPERATION_SETTINGS_NUMBERS_SENSOR_DATA_NUMBERS_BIT = 16

	SETTINGS_OPERATION_SETTINGS_TIMINGS_INTERRUPT_SLEEP_TIME_BIT = 0
	SETTINGS_OPERATION_SETTINGS_TIMINGS_INTERRUPT_OPEN_TIME_BIT  = 16

	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_GNSS_ENABLED_BIT                 = 0
	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_FLOWMETER_ENABLED_BIT            = 1
	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_TEMPERATURE_ENABLED_BIT          = 2
	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_HUMIDITY_ENABLED_BIT             = 3
	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_TEMPERATURE_HUMIDITY_ENABLED_BIT = 4
	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_MOTION_ENABLED_BIT               = 5
	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_SYSINFO_ENABLED_BIT              = 20
	SETTINGS_OPERATION_SETTINGS_SENSOR_DATA_EPOCH_ENABLED_BIT                = 21

	SETTINGS_LOGGING_SETTINGS_ENABLED = 0
)

func command_devices_subcommands_settings_subcommands_write_atnode(cCtx *cli.Context, filename string, output_to_stdout bool) (err error) {

	log.Printf("writing to file: %s\n", filename)

	/* generate sdcardsetting payload */
	settings, err := command_devices_subcommands_settings_subcommands_populate_settings("atnode")
	if err != nil {
		return
	}
	/**/

	/* prepare sdcard settings data */
	sdcard_settings := pb.SDCardSettings{
		OPERATION_SETTINGS_SETTINGS:             settings.get_settings().OperationSettings,
		OPERATION_SETTINGS_NUMBERS_SETTINGS:     settings.get_settings().OperationNumberSettings,
		OPERATION_SETTINGS_TIMINGS_SETTINGS:     settings.get_settings().OperationTimingsSettings,
		OPERATION_SETTINGS_SENSOR_DATA_SETTINGS: settings.get_settings().OperationSensorDataSettings,
		OPERATION_SETTINGS_LOGGING_SETTINGS:     settings.get_settings().OperationLoggingSettings,
	}

	/**/

	/* prepare protobuf data */
	sdcard_settings_bytes, err := proto.Marshal(&sdcard_settings)
	if err != nil {
		return err
	}
	/**/

	/* write the settings to file */
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	n, err := f.Write(sdcard_settings_bytes)
	if err != nil {
		return err
	}
	if n != len(sdcard_settings_bytes) {
		return fmt.Errorf("failed to write all bytes")
	}
	err = f.Close()
	if err != nil {
		return err
	}
	/**/
	fmt.Printf("settings succesfully written to %s\n", filename)
	return nil
}

func command_devices_subcommands_settings_subcommands_populate_settings(device_type string) (settings DeviceSettingsSurveyAnswerPayload, err error) {

	var qs = []*survey.Question{
		{
			Name: "OperationMode",
			Prompt: &survey.Select{
				Message: "Device Operation Mode:",
				Options: []string{"WebSocket", "Interrupt"},
				Default: settings.OperationMode,
			},
		},
		{
			Name: "LoggingEnable",
			Prompt: &survey.Confirm{
				Message: "Enable Logging: ",
				Default: false,
			},
		},
		{
			Name: "Sensors",
			Prompt: &survey.MultiSelect{
				Message: "Enable Sensors:",
				Options: []string{"GNSS", "Ambient Light", "Temperature", "Humidity", "Temperature/Humidity", "Inertial", "Magnetic", "Proximity", "Flowmeter", "Motion"},
				Default: []string{},
			},
		},
		{
			Name:   "InterruptOpenTime",
			Prompt: &survey.Input{Message: "Interrupt Open Time (Seconds):"},
		},

		{
			Name:   "InterruptSleepTime",
			Prompt: &survey.Input{Message: "Interrupt Sleep Time (Seconds):"},
		},
	}

	err = survey.Ask(qs, &settings)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}
