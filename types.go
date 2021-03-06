package main

import "time"

// Alert configuration in yaml
type DesiredAlertsConfig struct {
	Alerts []Alert `yaml:"alerts"`
}

// Alert configuration from api
type CurrentAlertsConfig struct {
	Results []ApiResult `json:"results"`
}

type ApiResult struct {
	Created time.Time `json:"created"`
	ID      string    `json:"id"`
	Updated time.Time `json:"updated"`
	Alert
}

type Alert struct {
	AlertName       string          `json:"alertName" yaml:"alertName"`
	Enabled         bool            `json:"enabled" yaml:"enabled"`
	MetricThreshold MetricThreshold `json:"metricThreshold" yaml:"metricThreshold"`
	Notifications   []Notifications `json:"notifications" yaml:"notifications"`
}

type MetricThreshold struct {
	MetricName string `json:"metricName" yaml:"metricName"`
	Operator   string `json:"operator" yaml:"operator"`
	Threshold  int    `json:"threshold" yaml:"threshold"`
	Units      string `json:"units" yaml:"units"`
}

type Notifications struct {
	NotificationType     string `json:"notificationType" yaml:"notificationType"`
	NotificationChannel  string `json:"notificationChannel,omitempty" yaml:"notificationChannel,omitempty"`
	DelayMin             int    `json:"delayMin" yaml:"delayMin"`
	IntervalMin          int    `json:"intervalMin" yaml:"intervalMin"`
	NotificationSchedule string `json:"notificationSchedule,omitempty" yaml:"notificationSchedule,omitempty"`
}

// The result of reconciling the desired and current configurations
type Action string

const (
	DELETE Action = "delete"
	CREATE Action = "create"
	UPDATE Action = "update"
)

type ReconcileActions struct {
	Actions []ReconcileAction `json:"reconcileActions"`
}

type ReconcileAction struct {
	AlertID string `json:"alertID"`
	Action  Action `json:"action"`
	Body    Alert  `json:"body"`
}
