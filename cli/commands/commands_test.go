package commands

import (
	"testing"

	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"github.com/spf13/cobra"
)

func TestAddCommands(t *testing.T) {
	cmd := &cobra.Command{}
	assert.False(t, cmd.HasSubCommands())
	AddCommands(cmd)
	assert.True(t, cmd.HasSubCommands())
	// Check if 'assessment_result' CMD is part of the sub commands
	for _, v := range cmd.Commands() {
		if v.Use == "assessment-result" {
			return
		}
	}
	t.Error("CMD 'assessment-result' is not part of the sub commands but should be")
}
