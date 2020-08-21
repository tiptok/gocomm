package eda

import "testing"

func TestExample(t *testing.T) {
	carSvr := &ServiceCar{}
	carSvr.Subscribe(&CarEngine{}, &CarFuelTank{})
	carSvr.Drive(50, 90)
	carSvr.Drive(100, 80)
	carSvr.Drive(120, 90)
}
