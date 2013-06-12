package virtualbox

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"time"
)

// This step starts the virtual machine.
//
// Uses:
//
// Produces:
type stepRun struct{
	vmName string
}

func (s *stepRun) Run(state map[string]interface{}) multistep.StepAction {
	driver := state["driver"].(Driver)
	ui := state["ui"].(packer.Ui)
	vmName := state["vmName"].(string)

	ui.Say("Starting the virtual machine...")
	command := []string{"startvm", vmName, "--type", "gui"}
	if err := driver.VBoxManage(command...); err != nil {
		ui.Error(fmt.Sprintf("Error starting VM: %s", err))
		return multistep.ActionHalt
	}

	s.vmName = vmName

	time.Sleep(15 * time.Second)
	return multistep.ActionContinue
}

func (s *stepRun) Cleanup(state map[string]interface{}) {
	if s.vmName == "" {
		return
	}

	driver := state["driver"].(Driver)
	ui := state["ui"].(packer.Ui)
	if err := driver.VBoxManage("controlvm", s.vmName, "poweroff"); err != nil {
		ui.Error(fmt.Sprintf("Error shutting down VM: %s", err))
	}
}