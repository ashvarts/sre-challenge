package main

import (
	"fmt"
	"math/rand"
	"time"
)

// for random alertId generation
func init() {
	rand.Seed(time.Now().UnixNano())
}

var charList = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

const IDLENGTH int = 24

func main() {
	fmt.Println("Hello, world.")
}

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

	missingAlerts := createMissingAlerts(desiredAlertsByName, currentAlertsByName)

	reconcileActions.Actions = append(reconcileActions.Actions, missingAlerts...)

	return reconcileActions
}

func createMissingAlerts(desiredAlertsByName map[string]Alert, currentAlertsByName map[string]ApiResult) []ReconcileAction {
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
