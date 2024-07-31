package bluetoothmanager

import (
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
)

type BluetoothManager struct {
	adapter   *adapter.Adapter1
	cancel    func()
	Discovery chan *adapter.DeviceDiscovered
}

func NewBluetoothManager() (*BluetoothManager, error) {
	adapter, err := adapter.GetDefaultAdapter()
	if err != nil {
		return nil, err
	}

	err = adapter.FlushDevices()
	if err != nil {
		return nil, err
	}

	return &BluetoothManager{
		adapter: adapter,
	}, nil
}

func (bm *BluetoothManager) GetBluetoothStatus() (bool, error) {
	props, err := bm.adapter.GetProperties()
	if err != nil {
		return false, err
	}

	return props.Powered, nil
}

func (bm *BluetoothManager) SetBluetoothPowerState(powered bool) error {
	return bm.adapter.SetProperty("Powered", powered)
}

func (bm *BluetoothManager) StartScan() error {
	discovery, cancel, err := api.Discover(bm.adapter, nil)
	if err != nil {
		return err
	}

	bm.Discovery = discovery
	bm.cancel = cancel

	return nil
}

func (bm *BluetoothManager) StopScan() {
	if bm.cancel != nil {
		bm.cancel()
	}
}
