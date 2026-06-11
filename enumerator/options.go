package enumerator

type Option func(o *enumerator)

func WithRolloverEnabled(r bool) Option {
	return func(o *enumerator) {
		o.rollover = r
	}
}
