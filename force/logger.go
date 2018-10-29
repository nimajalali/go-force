package force

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
)

func init() {

	// Initialises a global instance of a logrus logger.
	// Specific to the force package.
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
	}
	l.SetOutput(os.Stdout)
	logDefault := logrus.NewEntry(l)

	log = logDefault.WithFields(logrus.Fields{
		"Product": "ForceClient",
	})
}
