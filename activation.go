package neat

const (
	actMin = -100000.0
	actMax = 100000.0
)

func RandActivationFunc() *ActivationFunc {
	funcs = []func() *ActviationFunc{Linear, Sigmoid, Tanh, Sine, Gaussian}
	return funcs[rand.Intn(len(funcs))]()
}

type ActivationFunc struct {
	name string
	f    func(x float64) float64
}

func Linear() *ActivationFunc {
	return &ActivationFunc{
		name: "linear",
		f: func(x float64) float64 {
			return x
		},
	}
}

func Sigmoid() *ActivationFunc {
	return &ActivationFunc{
		name: "sigmoid",
		f: func(x float64) float64 {
			return 1.0 / (1.0 + math.Exp(-x))
		},
	}
}

func Tanh() *ActivationFunc {
	return &ActivationFunc{
		name: "tanh",
		f: func(x float64) float64 {
			return math.Tanh(x)
		},
	}
}

func Sine() *ActivationFunc {
	return &ActivationFunc{
		name: "sine",
		f: func(x float64) float64 {
			return math.Sin(x)
		},
	}
}

func Gaussian() *ActivationFunc {
	return &ActivationFunc{
		name: "gaussian",
		f: func(x float64) float64 {
			return math.Exp(-x * x)
		},
	}
}

func (a *ActivationFunc) Name() float64 {
	return a.name
}

func (a *ActivationFunc) Activate(x float64) float64 {
	clipped = math.Min(math.Max(x, actMin), actMax)
	return a.f(x)
}
