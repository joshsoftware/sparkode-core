package helpers

type Status struct {
	Id    int
	Value string
}

var ErrorToStatus = map[string]Status{
	"queue":   Status{Id: 1, Value: "In Queue"},
	"process": Status{Id: 2, Value: "Processing"},
	"ac":      Status{Id: 3, Value: "Accepted"},
	"wa":      Status{Id: 4, Value: "Wrong Answer"},
	"tle":     Status{Id: 5, Value: "Time Limit Exceeded"},
	"ce":      Status{Id: 6, Value: "Compilation Error"},
	"sigsegv": Status{Id: 7, Value: "Runtime Error (SIGSEGV)"},
	"sigxfsz": Status{Id: 8, Value: "Runtime Error (SIGXFSZ)"},
	"sigfpe":  Status{Id: 9, Value: "Runtime Error (SIGFPE)"},
	"sigabrt": Status{Id: 10, Value: "Runtime Error (SIGABRT)"},
	"nzec":    Status{Id: 11, Value: "Runtime Error (NZEC)"},
	"other":   Status{Id: 12, Value: "Runtime Error (Other)"},
	"boxerr":  Status{Id: 13, Value: "Internal Error"},
	"exeerr":  Status{Id: 14, Value: "Exec Format Error"},
}

var StatusIDToError = map[int]string{
	1:  "queue",
	2:  "process",
	3:  "ac",
	4:  "wa",
	5:  "tle",
	6:  "ce",
	7:  "sigsegv",
	8:  "sigxfsz",
	9:  "sigfpe",
	10: "sigabrt",
	11: "nzec",
	12: "other",
	13: "boxerr",
	14: "exeerr",
}

func GetErrorByStatusCode(id int) (status Status) {
	err, ok := StatusIDToError[id]
	if ok {
		status, ok = ErrorToStatus[err]
	}
	return status
}
