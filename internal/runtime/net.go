package runtime

import "fmt"

func GetAdvertiseAddress() string {
	return fmt.Sprintf(
		"%s:%d",
		Conf.Spec.Advertise,
		Conf.Spec.Port,
	)
}
