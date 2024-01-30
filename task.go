package workerpool

import "context"

type Task func(context.Context)
