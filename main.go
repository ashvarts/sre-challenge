package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"

	"gopkg.in/yaml.v2"
)

const IDLENGTH int = 24

var charList = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var desiredConfig DesiredAlertsConfig
	var currentConfig CurrentAlertsConfig

	currentConfigFilePath := flag.String("current-config", "", "required: the path to the current config file (json api result)")
	desiredConfigFilePath := flag.String("desired-config", "", "required: the path to the desired config file (yaml configuration)")
	flag.Parse()

	if (*currentConfigFilePath == "") || (*desiredConfigFilePath == "") {
		flag.PrintDefaults()
		os.Exit(0)
	}

	currentConfigFile, err := ioutil.ReadFile(*currentConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(currentConfigFile, &currentConfig)
	if err != nil {
		log.Fatal("something went wrong reading current-config (json)", err)
	}

	desiredConfigFile, err := ioutil.ReadFile(*desiredConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(desiredConfigFile, &desiredConfig)
	if err != nil {
		log.Fatal("something went wrong reading desired-config (yaml)", err)
	}

	reconcileActions := Reconcile(desiredConfig, currentConfig)
	actions, err := json.MarshalIndent(reconcileActions.Actions, "", " ")
	if err != nil {
		log.Fatal("something went wrong with generating reconcile plan", err)
	}
	summary := Summary(reconcileActions)

	fmt.Printf("%s\n", actions)
	fmt.Printf("Summary: {Created:%d,Deleted:%d,Updated:%d}\n", summary[CREATE], summary[DELETE], summary[UPDATE])

}

func Summary(actions ReconcileActions) map[Action]int { //TODO: Test
	summary := make(map[Action]int)
	for _, action := range actions.Actions {
		summary[action.Action]++
	}
	return summary
}

// Reconcile takes a desired and current configuration and returns a list of actions needed
// to match the current configuration to the desired configuration.
func Reconcile(desired DesiredAlertsConfig, current CurrentAlertsConfig) ReconcileActions {
	reconcileActions := ReconcileActions{}

	desiredAlerts := desired.Alerts
	currentAlerts := current.Results

	// if current alerts are empty, create all alerts in desired config
	if len(currentAlerts) == 0 {
		for _, alert := range desiredAlerts {
			alertId := RandId(IDLENGTH)
			action := CreateReconcileAction(alertId, CREATE, alert)
			reconcileActions.Actions = append(reconcileActions.Actions, action)
		}
		return reconcileActions
	}

	desiredAlertsByName := ConfigByAlertName(desiredAlerts)
	currentAlertsByName := ResultsByAlertName(currentAlerts)

	missingAlerts := createActionsForMissingAlerts(desiredAlertsByName, currentAlertsByName)
	if len(missingAlerts) > 0 {
		reconcileActions.Actions = append(reconcileActions.Actions, missingAlerts...)
	}

	changedAlerts := createActionsForUpdatedAlerts(desiredAlertsByName, currentAlertsByName)
	if len(changedAlerts) > 0 {
		reconcileActions.Actions = append(reconcileActions.Actions, changedAlerts...)
	}

	deletedAlerts := createActionsForDeletedAlerts(desiredAlertsByName, currentAlertsByName)
	if len(deletedAlerts) > 0 {
		reconcileActions.Actions = append(reconcileActions.Actions, deletedAlerts...)
	}
	return reconcileActions
}
func createActionsForDeletedAlerts(desiredAlertsByName map[string]Alert, currentAlertsByName map[string]ApiResult) []ReconcileAction {
	var reconcileActions []ReconcileAction
	for alertName, currentAlert := range currentAlertsByName {
		if _, ok := desiredAlertsByName[alertName]; !ok {
			alertID := currentAlert.ID
			action := CreateReconcileAction(alertID, DELETE, currentAlert.Alert)
			reconcileActions = append(reconcileActions, action)
		}
	}
	return reconcileActions
}

func createActionsForUpdatedAlerts(desiredAlertsByName map[string]Alert, currentAlertsByName map[string]ApiResult) []ReconcileAction {
	var reconcileActions []ReconcileAction
	for alertName, desiredAlert := range desiredAlertsByName {
		if currentAlert, ok := currentAlertsByName[alertName]; ok {
			if !reflect.DeepEqual(currentAlert.Alert, desiredAlert) {
				alertID := currentAlert.ID
				action := CreateReconcileAction(alertID, UPDATE, desiredAlert)
				reconcileActions = append(reconcileActions, action)
			}
		}
	}
	return reconcileActions
}

func createActionsForMissingAlerts(desiredAlertsByName map[string]Alert, currentAlertsByName map[string]ApiResult) []ReconcileAction {
	var reconcileActions []ReconcileAction
	for alertName, alert := range desiredAlertsByName {
		if _, ok := currentAlertsByName[alertName]; !ok {
			alertId := RandId(IDLENGTH)
			action := CreateReconcileAction(alertId, CREATE, alert)
			reconcileActions = append(reconcileActions, action)
		}
	}
	return reconcileActions
}

func CreateReconcileAction(alertId string, action Action, body Alert) ReconcileAction {
	reconcileAction := ReconcileAction{
		AlertID: alertId,
		Action:  action,
		Body:    body,
	}

	return reconcileAction
}

func ConfigByAlertName(alerts []Alert) map[string]Alert {
	alertsByName := make(map[string]Alert, len(alerts))
	for _, alert := range alerts {
		alertsByName[alert.AlertName] = alert
	}

	return alertsByName
}

func ResultsByAlertName(results []ApiResult) map[string]ApiResult {
	resultsByName := make(map[string]ApiResult, len(results))
	for _, result := range results {
		resultsByName[result.AlertName] = result
	}

	return resultsByName
}

// generate a random id of length n
func RandId(n int) string {
	id := make([]rune, n)
	for i := range id {
		id[i] = charList[rand.Intn(len(charList))]
	}
	return string(id)
}
