package enumerator

type Option func(o *enumerator)

var (
	defaultOptions = []Option{
		WithRolloverEnabled(false),
	}
)

func WithRolloverEnabled(r bool) Option {
	return func(o *enumerator) {
		o.rollover = r
	}
}
