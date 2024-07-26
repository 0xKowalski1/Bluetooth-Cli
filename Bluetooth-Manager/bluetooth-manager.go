package bluetoothmanager

import (
        "github.com/muka/go-bluetooth/bluez/profile/adapter"
)

type BluetoothManager struct {
        adapter *adapter.Adapter1
}

func NewBluetoothManager() (*BluetoothManager, error){
        adapter, err := adapter.GetDefaultAdapter()
        if err != nil {
                return nil, err
        }
        return &BluetoothManager{
                adapter:adapter,
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

