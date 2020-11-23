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
