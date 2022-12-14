package events

import (
	"fmt"
	"os"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

func Run(sm *session.SessionManager) {
	mqtt.Subscribe("mega/state-report", func(topic string, payload []byte) {
		t := time.Now().UnixMicro()

		sr := &machinepb.StateReport{}
		err := proto.Unmarshal(payload, sr)
		if err != nil {
			fmt.Printf("error unmarshalling state report: %v\n", err)
			return
		}
		sr.TimestampUnixMicros = uint64(t)
		fmt.Printf("%+v\n", sr)

		// Abort unless session is active
		session, _ := sm.GetLatestSession()
		if session == nil || session.Complete || session.Paused {
			return
		}

		saveSessionStateReport(session, sr)
	})
}

func saveSessionStateReport(s *session.Session, sr *machinepb.StateReport) {
	list := []*machinepb.StateReport{sr}
	output, err := yaml.Marshal(list)
	if err != nil {
		fmt.Printf("error marshalling state report to yaml: %v\n", err)
	}

	p := filesystem.GetStateReportPath(uint64(s.Id))

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file for state report storage: %v\n", err)
	}
	defer f.Close()
	f.Write(output)
}
