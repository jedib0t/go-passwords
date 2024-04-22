package odometer

type Option func(o *odometer)

var (
	defaultOptions = []Option{
		WithRolloverEnabled(false),
	}
)

func WithRolloverEnabled(r bool) Option {
	return func(o *odometer) {
		o.rollover = r
	}
}
