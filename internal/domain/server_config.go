package domain

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"tailscale.com/tailcfg"
	"tailscale.com/types/key"
)

type configKey string

const (
	controlKeysConfigKey configKey = "control_keys"
	derpMapConfigKey     configKey = "derp_map"
)

type ControlKeys struct {
	ControlKey       key.MachinePrivate
	LegacyControlKey key.MachinePrivate
}

type ServerConfig struct {
	Key   configKey `gorm:"primary_key"`
	Value []byte
}

func (r *repository) GetControlKeys(ctx context.Context) (*ControlKeys, error) {

	var m ControlKeys
	err := r.getServerConfig(ctx, controlKeysConfigKey, &m)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *repository) SetControlKeys(ctx context.Context, v *ControlKeys) error {
	return r.setServerConfig(ctx, controlKeysConfigKey, v)
}

func (r *repository) GetDERPMap(ctx context.Context) (*tailcfg.DERPMap, error) {
	var m tailcfg.DERPMap

	err := r.getServerConfig(ctx, derpMapConfigKey, &m)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *repository) SetDERPMap(ctx context.Context, v *tailcfg.DERPMap) error {
	return r.setServerConfig(ctx, "derp_map", v)
}

func (r *repository) getServerConfig(ctx context.Context, s configKey, v interface{}) error {
	var m ServerConfig
	tx := r.withContext(ctx).Take(&m, "key = ?", s)

	if tx.Error != nil {
		return tx.Error
	}

	err := json.Unmarshal(m.Value, v)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) setServerConfig(ctx context.Context, s configKey, v interface{}) error {
	marshal, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c := &ServerConfig{
		Key:   s,
		Value: marshal,
	}
	tx := r.withContext(ctx).Save(c)

	return tx.Error
}
