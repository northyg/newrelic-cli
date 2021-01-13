package execution

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type TerminalStatusReporter struct {
	status *StatusRollup
}

// NewTerminalStatusReporter is an implementation of the ExecutionStatusReporter interface that reports execution status to STDOUT.
func NewTerminalStatusReporter() *TerminalStatusReporter {
	rollup := NewStatusRollup()
	rollup.LogFilePath = config.DefaultConfigDirectory + "/" + config.DefaultLogFile
	r := TerminalStatusReporter{
		status: &rollup,
	}

	return &r
}

func (r TerminalStatusReporter) ReportRecipeFailed(event RecipeStatusEvent) error {
	r.status.withRecipeEvent(event, StatusTypes.FAILED)
	return nil
}

func (r TerminalStatusReporter) ReportRecipeInstalled(event RecipeStatusEvent) error {
	r.status.withRecipeEvent(event, StatusTypes.INSTALLED)
	return nil
}

func (r TerminalStatusReporter) ReportRecipeSkipped(event RecipeStatusEvent) error {
	r.status.withRecipeEvent(event, StatusTypes.SKIPPED)
	return nil
}

func (r TerminalStatusReporter) ReportRecipesAvailable(recipes []types.Recipe) error {
	r.status.withAvailableRecipes(recipes)
	return nil
}

func (r TerminalStatusReporter) ReportRecipeAvailable(recipe types.Recipe) error {
	r.status.withAvailableRecipe(recipe)
	return nil
}

func (r TerminalStatusReporter) ReportComplete() error {

	if r.hasFailed() {
		return fmt.Errorf("one or more integrations failed to install, check the install log for more details: %s", r.status.LogFilePath)
	}

	msg := `
	Success! Your data is available in New Relic.

	Go to New Relic to confirm and start exploring your data.`

	fmt.Println(msg)

	for _, entityGUID := range r.status.EntityGUIDs {
		fmt.Printf("\n\thttps://one.newrelic.com/redirect/entity/%s\n", entityGUID)
	}

	return nil
}

func (r TerminalStatusReporter) hasFailed() bool {
	for _, s := range r.status.Statuses {
		if s.Status == StatusTypes.FAILED {
			return true
		}
	}

	return false
}