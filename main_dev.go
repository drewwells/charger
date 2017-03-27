// build +dev

package main

func init() {
	parseCred = func(c config) Cred {
		return c.Test
	}
}
