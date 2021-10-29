package runtime

import (
	"github.com/ahmetsoykan/controllers-tutorials/001/pkg/subscription"
)

func RunLoop(subscriptions []subscription.ISubscription) error {

	for _, subscription := range subscriptions {

		wiface, err := subscription.Subscribe()
		if err != nil {
			return err
		}

		go func() {

			for {
				select {
				case msg := <-wiface.ResultChan():
					subscription.Reconcile(msg.Object, msg.Type)
				}
			}

		}() // Could signal handler into them??
	}

	for _, subscription := range subscriptions {

		select {
		case _ = <-subscription.IsComplete():
			break
		}
	}
	return nil
}
