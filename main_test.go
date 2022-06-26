package main

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReconcileAction(t *testing.T) {
	id := RandId(24)
	testAlert := newTestAlert()
	type args struct {
		alertId string
		action  Action
		body    Alert
	}
	tests := []struct {
		name string
		args args
		want ReconcileAction
	}{
		{
			name: "Create",
			args: args{
				alertId: id,
				action:  CREATE,
				body:    testAlert,
			},
			want: ReconcileAction{
				AlertID: id,
				Action:  CREATE,
				Body:    testAlert,
			},
		},
		{
			name: "Update",
			args: args{
				alertId: id,
				action:  UPDATE,
				body:    testAlert,
			},
			want: ReconcileAction{
				AlertID: id,
				Action:  UPDATE,
				Body:    testAlert,
			},
		},
		{
			name: "Delete",
			args: args{
				alertId: id,
				action:  DELETE,
				body:    testAlert,
			},
			want: ReconcileAction{
				AlertID: id,
				Action:  DELETE,
				Body:    testAlert,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateReconcileAction(tt.args.alertId, tt.args.action, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateReconcileAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReconcile(t *testing.T) {
	testAlert0 := newTestAlert()
	testAlert1 := newTestAlert()
	testAlert2 := newTestAlert()
	testAlert3 := newTestAlert()

	t.Run("Reconcile should create 'create' actions when current config is empty", func(t *testing.T) {
		desiredAlerts := DesiredAlertsConfig{
			Alerts: []Alert{
				testAlert0,
				testAlert1,
				testAlert2,
				testAlert3,
			},
		}
		currentConfig := CurrentAlertsConfig{}

		reconcileActions := Reconcile(desiredAlerts, currentConfig)

		if len(reconcileActions.Actions) != len(desiredAlerts.Alerts) {
			t.Errorf("reconcile actions count should be: %d, got: %d", len(desiredAlerts.Alerts), len(reconcileActions.Actions))
		}

		for idx, action := range reconcileActions.Actions {
			if action.AlertID == "" {
				t.Errorf("alertId should not be empty")
			}

			if action.Action != CREATE {
				t.Errorf("action should be 'create', got: %s", action.Action)
			}

			if !reflect.DeepEqual(action.Body, desiredAlerts.Alerts[idx]) {
				t.Errorf("action body should equal: %+v, got: %+v", desiredAlerts.Alerts[idx], action.Body)
			}
		}

	})

	t.Run("Reconcile should only create 'create' actions if current config doesn't exists", func(t *testing.T) {
		desiredAlerts := DesiredAlertsConfig{
			Alerts: []Alert{
				testAlert0,
				testAlert1,
				testAlert2,
				testAlert3,
			},
		}
		currentConfig := newCurrentAlertsConfig(
			[]Alert{
				testAlert0,
				testAlert1,
			})

		reconcileActions := Reconcile(desiredAlerts, currentConfig)
		if len(reconcileActions.Actions) != 2 {
			t.Errorf("reconcile actions count should be: %d, got: %d", 2, len(reconcileActions.Actions))
		}

		for _, action := range reconcileActions.Actions {
			if action.AlertID == "" {
				t.Errorf("alertId should not be empty")
			}
			if action.Action != CREATE {
				t.Errorf("action should be 'create', got: %s", action.Action)
			}
		}
		assertAlertDefinitionMatchesResultBody(t, reconcileActions, []Alert{testAlert2, testAlert3})

	})

	t.Run("Reconcile should create 'update' action for changed ", func(t *testing.T) {
		currentConfig := newCurrentAlertsConfig(
			[]Alert{
				testAlert0,
				testAlert1,
			})

		testAlert0.MetricThreshold.Operator = "updated_operator"
		desiredAlerts := DesiredAlertsConfig{
			Alerts: []Alert{
				testAlert0,
				testAlert1,
			},
		}

		reconcileActions := Reconcile(desiredAlerts, currentConfig)

		if len(reconcileActions.Actions) != 1 {
			t.Errorf("reconcile actions count should be: %d, got: %d", 1, len(reconcileActions.Actions))
		}

		for _, action := range reconcileActions.Actions {
			if action.AlertID != currentConfig.Results[0].ID {
				t.Errorf("reconcile action alert ID should match id in result. want: %s, got: %s", currentConfig.Results[0].ID, action.AlertID)
			}
			if action.Action != UPDATE {
				t.Errorf("action should be 'create', got: %s", action.Action)
			}
		}

		assertAlertDefinitionMatchesResultBody(t, reconcileActions, []Alert{testAlert0})

	})

	t.Run("Reconcile should create 'delete' action for removed desired alert", func(t *testing.T) {
		idxOfRemovedAlert := 1
		currentConfig := newCurrentAlertsConfig(
			[]Alert{
				testAlert0,
				testAlert1,
				testAlert2,
			})

		desiredAlerts := DesiredAlertsConfig{
			Alerts: []Alert{
				testAlert0,
				testAlert2,
			},
		}

		reconcileActions := Reconcile(desiredAlerts, currentConfig)

		if len(reconcileActions.Actions) != 1 {
			t.Errorf("reconcile actions count should be: %d, got: %d", 1, len(reconcileActions.Actions))
		}

		for _, action := range reconcileActions.Actions {
			if action.AlertID != currentConfig.Results[idxOfRemovedAlert].ID {
				t.Errorf("reconcile action alert ID should match id in result. want: %s, got: %s", currentConfig.Results[idxOfRemovedAlert].ID, action.AlertID)
			}
			if action.Action != DELETE {
				t.Errorf("action should be 'create', got: %s", action.Action)
			}
		}

		assertAlertDefinitionMatchesResultBody(t, reconcileActions, []Alert{currentConfig.Results[idxOfRemovedAlert].Alert})

	})

}

func assertAlertDefinitionMatchesResultBody(t *testing.T, reconcileActions ReconcileActions, alerts []Alert) {
	var reconcileAlertBodies []Alert
	for _, v := range reconcileActions.Actions {
		reconcileAlertBodies = append(reconcileAlertBodies, v.Body)
	}
	assert.ElementsMatchf(t, reconcileAlertBodies, alerts, "reconcile alert bodies should match desired alerts")

}

// test helper
func newTestAlert() Alert {
	randSuffix := RandId(5)
	return Alert{
		AlertName: "testAlertName" + randSuffix,
		Enabled:   true,
		MetricThreshold: MetricThreshold{
			MetricName: "testMetricName" + randSuffix,
			Operator:   "testOperator" + randSuffix,
			Threshold:  rand.Intn(60),
			Units:      "testUnits" + randSuffix,
		},
		Notifications: []Notifications{
			{
				NotificationType:     "testNotificationType" + randSuffix,
				NotificationChannel:  "testNotificationChannel" + randSuffix,
				DelayMin:             rand.Intn(60),
				IntervalMin:          rand.Intn(60),
				NotificationSchedule: "testNotificationSchedule" + randSuffix,
			},
		},
	}
}

func newCurrentAlertsConfig(alerts []Alert) CurrentAlertsConfig {
	currentResultsConfig := CurrentAlertsConfig{}
	for _, alert := range alerts {
		id := RandId(IDLENGTH)
		apiResult := ApiResult{}
		apiResult.ID = id
		apiResult.Alert = alert
		currentResultsConfig.Results = append(currentResultsConfig.Results, apiResult)
	}

	return currentResultsConfig
}
