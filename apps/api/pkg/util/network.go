package util

type Network string

const (
	Mainnet        Network = "mainnet"
	CalibrationNet Network = "calibration"
)

func (n Network) IsMainnet() bool {
	return n == Mainnet
}

func (n Network) String() string {
	return string(n)
}

func StrToNetwork(network string) Network {
	if network == Mainnet.String() {
		return Mainnet
	}

	return CalibrationNet
}
