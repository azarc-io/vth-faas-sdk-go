package sparkv1

import "time"

const maxInactiveConsumerDuration = time.Hour
const maxInactiveResetConsumerDuration = maxInactiveConsumerDuration / 2
const maxConsumerFetchWait = time.Second * 15
