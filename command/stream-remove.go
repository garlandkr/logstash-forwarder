package command

import (
	"fmt"
	"lsf"
	"lsf/anomaly"
	"lsf/schema"
	"lsf/system"
	"os"
)

const removeStreamCmdCode lsf.CommandCode = "stream-remove"

type removeStreamOptionsSpec struct {
	global BoolOptionSpec
	id     StringOptionSpec
}

var removeStream *lsf.Command
var removeStreamOptions *removeStreamOptionsSpec

func init() {

	removeStream = &lsf.Command{
		Name:  removeStreamCmdCode,
		About: "Remove a new log stream",
		Init:  verifyRemoveStreamRequiredOpts,
		Run:   runRemoveStream,
		Flag:  FlagSet(removeStreamCmdCode),
	}
	removeStreamOptions = &removeStreamOptionsSpec{
		global: NewBoolFlag(removeStream.Flag, "G", "global", false, "global scope operation", false),
		id:     NewStringFlag(removeStream.Flag, "s", "stream-id", "", "unique identifier for stream", true),
	}
}
func verifyRemoveStreamRequiredOpts(env *lsf.Environment, args ...string) error {
	if e := verifyRequiredOption(removeStreamOptions.id); e != nil {
		return e
	}
	return nil
}

// REVU: TODO definitively require a stream 'x' lock for use by
// processes that expect the stream (info) to remain in place.
// For now, assuming this is the same "stream.<name>.stream.lock"
// lock file.
func runRemoveStream(env *lsf.Environment, args ...string) (err error) {
	anomaly.Recover(&err)

	id := schema.StreamId(*removeStreamOptions.id.value)

	// check existing
	docid := system.DocId(fmt.Sprintf("stream.%s.stream", id))
	doc, e := env.LoadDocument(docid)
	if e != nil || doc == nil {
		return lsf.E_NOTEXISTING
	}

	// lock lsf port's "streams" resource
	lockid := env.ResourceId("streams")
	lock, ok, e := system.LockResource(lockid, "add stream "+string(id))
	anomaly.PanicOnError(e, "command.runRemoveStream:", "lockResource:")
	anomaly.PanicOnFalse(ok, "command.runRemoveStream:", "lockResource:", string(id))
	defer lock.Unlock()

	// remove doc
	ok, e = env.DeleteDocument(docid)
	anomaly.PanicOnError(e, "command.runRemoveStream:", "DeleteDocument:", string(id))
	anomaly.PanicOnFalse(ok, "command.runRemoveStream:", "DeleteDocument:", string(id))

	// remove the stream's directory
	// REVU: this command needs a check to see if any procs
	// related to this stream are running . OK for initial.
	dir, fname := system.DocpathForKey(env.Port(), docid)
	fmt.Printf("DEBUG: runRemoveStream: %s %s\n", dir, fname)

	e = os.RemoveAll(dir)
	anomaly.PanicOnError(e, "command.runRemoveStream:", "os.RemoveAll:", dir)

	return nil
}