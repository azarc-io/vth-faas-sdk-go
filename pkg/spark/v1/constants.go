package sparkv1

import "time"

const maxInactiveConsumerDuration = time.Hour
const maxInactiveResetConsumerDuration = maxInactiveConsumerDuration / 2
const maxConsumerFetchWait = time.Second * 15
const ConsumerBatch = 15
const maxConsumerDeliver = 1
const maxConsumerAckPending = 1
const maxConsumerCreationRetries = 3
