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
	"fmt"
	"github.com/bozzo/knx-go/knx"
	"github.com/bozzo/knx-go/knx/cemi"
	"github.com/bozzo/knx-go/knx/dpt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"reflect"
)

var floatType = reflect.TypeOf(float64(0))
var boolType = reflect.TypeOf(false)

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if v.Type().ConvertibleTo(boolType) {
		if v.Convert(boolType).Bool() {
			return 1.0, nil
		}
		return 0.0, nil
	}
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func getDptByGroupAddr(dest cemi.GroupAddr) (string, error) {
	dptID := config.getDptID(dest.String())
	if dptID == "" {
		return "", fmt.Errorf("type for %s not found", dest.String())
	}
	return dptID, nil
}

func parseLDataInd(event knx.GroupEvent) (object dpt.DatapointValue, err error) {
	dptID, err := getDptByGroupAddr(event.Destination)
	if err != nil {
		logrus.Errorf("Type for %s not found", event.Destination)
		return nil, err
	}

	object, found := dpt.Produce(dptID)
	if !found {
		logrus.Errorf("Type for %s not found", dptID)
		return nil, fmt.Errorf("type for %s not found", dptID)
	}
	logrus.Debugf("Produced object: %+v (%s)", object, dptID)

	err = object.Unpack(event.Data)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Unpacked value: %+v)", object)

	value, err := getFloat(object)
	if err != nil {
		logrus.Errorf("Object can't be converted to float64: %v", err)
		return nil, err
	}
	logrus.Debugf("Float value: %f (original: %+v)", value, object)

	getMetric(config.getDptName(dptID), dptID).With(prometheus.Labels{"source": event.Source.String(), "destination": event.Destination.String()}).Set(value)

	return object, nil
}
