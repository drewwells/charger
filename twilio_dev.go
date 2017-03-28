// build +dev

package charger

func init() {
	parseCred = func(c config) Cred {
		return c.Test
	}
}
