/*
 *
 *    Copyright 2020 Boris Barnier <bozzo@users.noreply.github.com>
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"github.com/bozzo/knx-go/knx/cemi"
	"github.com/bozzo/knx-go/knx/dpt"
	"testing"
)

func Test_getFloat(t *testing.T) {
	//type args struct {
	//	unk interface{}
	//}
	tests := []struct {
		name    string
		args    interface{}
		want    float64
		wantErr bool
	}{
		{"false should be transformed into float", dpt.DPT_1001(false), 0.0, false},
		{"true should be transformed into float", dpt.DPT_1001(true), 1.0, false},
		{"float should be transformed into float", dpt.DPT_14001(12.00), 12.0, false},
		{"int should be transformed into float", dpt.DPT_12001(12), 12.0, false},
		{"string should be raise an error", "string is not a float", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFloat(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDptByGroupAddr(t *testing.T) {
	tests := []struct {
		name    string
		args    cemi.GroupAddr
		want    string
		wantErr bool
	}{
		{"Group 0/0/1 should return DPT 1.001", cemi.GroupAddr(1), "1.001", false},
		{"Group 0/0/2 should return DPT 12.001", cemi.GroupAddr(2), "12.001", false},
		{"Group 0/0/3 should return DPT 14.001", cemi.GroupAddr(3), "14.001", false},
		{"Unknown group 0/6/1 should return an error", cemi.GroupAddr(0x601), "", true},
	}

	config.parseConfig(YmlConfiguration{
		Version: 1,
		Mapping: []YmlConfigurationTypeMapping{
			{"1.001", []string{"0/0/1"}, "switch"},
			{"12.001", []string{"0/0/2"}, "counter"},
			{"14.001", []string{"0/0/3"}, "Acceleration Angular"},
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDptByGroupAddr(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDptByGroupAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDptByGroupAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}