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
