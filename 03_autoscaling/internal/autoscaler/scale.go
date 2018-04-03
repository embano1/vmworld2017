package autoscaler

import (
	"fmt"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (as *AutoScaler) scale(qcurrent int) error {
	var currentRepl, newRepl int32

	if qcurrent < upper && qcurrent > lower {
		as.log.Printf("Queue %q within thresholds (current: %d, lower: %d, upper: %d)", queuename, qcurrent, lower, upper)
		return nil
	}

	di := as.kclient.Apps().Deployments(as.ns)
	d, err := di.Get(as.dpl, meta_v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not get deployment %q: %v", as.dpl, err)
	}
	currentRepl = *d.Spec.Replicas

	// scale up
	if qcurrent >= upper {
		as.log.Printf("Queue %q above threshold (current: %d, upper: %d), scaling up deployment %q", queuename, qcurrent, upper, as.dpl)
		if int(currentRepl) > as.w {
			as.log.Printf("Max workers (%d) already running, skipping.", as.w)
			return nil
		}

		newRepl = currentRepl + 1
		d.Spec.Replicas = &newRepl
		_, err = di.Update(d)
		if err != nil {
			as.log.Printf("could not update deployment %q: %v", as.dpl, err)
		}
	}

	// scale down
	if qcurrent <= lower {
		if currentRepl == 0 {
			as.log.Printf("Queue %q below threshold (current: %d, lower: %d) and no workers running...skipping", queuename, qcurrent, lower)
			return nil
		}
		as.log.Printf("Queue %q below threshold (current: %d, lower: %d), scaling down deployment %q", queuename, qcurrent, lower, as.dpl)

		newRepl = currentRepl - 1
		d.Spec.Replicas = &newRepl
		_, err = di.Update(d)
		if err != nil {
			as.log.Printf("could not update deployment %q: %v", as.dpl, err)
		}
	}

	return nil
}
