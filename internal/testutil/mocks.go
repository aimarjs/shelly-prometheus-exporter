package testutil

import (
	"github.com/aimar/shelly-prometheus-exporter/internal/client"
)

// CreateMockStatusResponse creates a mock StatusResponse for testing
func CreateMockStatusResponse(mac string, uptime int, ramSize, ramFree, fsSize, fsFree int, power, energy float64, temp float64) client.StatusResponse {
	return client.StatusResponse{
		Sys: struct {
			Mac              string `json:"mac"`
			RestartRequired  bool   `json:"restart_required"`
			Time             string `json:"time"`
			Unixtime         int64  `json:"unixtime"`
			LastSyncTs       int64  `json:"last_sync_ts"`
			Uptime           int    `json:"uptime"`
			RAMSize          int    `json:"ram_size"`
			RAMFree          int    `json:"ram_free"`
			RAMMinFree       int    `json:"ram_min_free"`
			FSSize           int    `json:"fs_size"`
			FSFree           int    `json:"fs_free"`
			CfgRev           int    `json:"cfg_rev"`
			KvsRev           int    `json:"kvs_rev"`
			ScheduleRev      int    `json:"schedule_rev"`
			WebhookRev       int    `json:"webhook_rev"`
			BtrelayRev       int    `json:"btrelay_rev"`
			AvailableUpdates struct {
				Stable struct {
					Version string `json:"version"`
				} `json:"stable"`
			} `json:"available_updates"`
			ResetReason int `json:"reset_reason"`
		}{
			Mac:     mac,
			Uptime:  uptime,
			RAMSize: ramSize,
			RAMFree: ramFree,
			FSSize:  fsSize,
			FSFree:  fsFree,
		},
		Temperature: struct {
			ID int     `json:"id"`
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		}{
			TC: temp,
		},
		EM: struct {
			ID             int      `json:"id"`
			ACurrent       float64  `json:"a_current"`
			AVoltage       float64  `json:"a_voltage"`
			AActPower      float64  `json:"a_act_power"`
			AAprtPower     float64  `json:"a_aprt_power"`
			APF            float64  `json:"a_pf"`
			AFreq          float64  `json:"a_freq"`
			BCurrent       float64  `json:"b_current"`
			BVoltage       float64  `json:"b_voltage"`
			BActPower      float64  `json:"b_act_power"`
			BAprtPower     float64  `json:"b_aprt_power"`
			BPF            float64  `json:"b_pf"`
			BFreq          float64  `json:"b_freq"`
			CCurrent       float64  `json:"c_current"`
			CVoltage       float64  `json:"c_voltage"`
			CActPower      float64  `json:"c_act_power"`
			CAprtPower     float64  `json:"c_aprt_power"`
			CPF            float64  `json:"c_pf"`
			CFreq          float64  `json:"c_freq"`
			NCurrent       *float64 `json:"n_current"`
			TotalCurrent   float64  `json:"total_current"`
			TotalActPower  float64  `json:"total_act_power"`
			TotalAprtPower float64  `json:"total_aprt_power"`
		}{
			TotalActPower: power,
		},
		EMData: struct {
			ID                 int     `json:"id"`
			ATotalActEnergy    float64 `json:"a_total_act_energy"`
			ATotalActRetEnergy float64 `json:"a_total_act_ret_energy"`
			BTotalActEnergy    float64 `json:"b_total_act_energy"`
			BTotalActRetEnergy float64 `json:"b_total_act_ret_energy"`
			CTotalActEnergy    float64 `json:"c_total_act_energy"`
			CTotalActRetEnergy float64 `json:"c_total_act_ret_energy"`
			TotalAct           float64 `json:"total_act"`
			TotalActRet        float64 `json:"total_act_ret"`
		}{
			TotalAct: energy,
		},
	}
}