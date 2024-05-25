package gonweb

type FilterBuilder func(next GonHandlerFunc) GonHandlerFunc
