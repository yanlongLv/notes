package main

type Driver struct {
	d []string
	n string
}

func main() {
	aa := make([]string, 10)
	aa[0] = "a"
	driver := &Driver{}
	driver.setTrip(aa)
}

func (d *Driver) setTrip(dd []string) {
	d.d = dd
}
