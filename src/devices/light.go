package devices

import "hub"

type LightAttributes struct {
	device_attributes     Attributes
	startup_on_off        string
	is_on                 bool
	light_Level           int
	color_temperature     int
	color_temperature_min int
	color_temperature_max int
	color_hue             float32
	color_saturation      float32
}

type Light struct {
	device          Device
	dirigera_client hub.Hub
	attributes      LightAttributes
}

func reload(l Light) Light {
  data := l.dirigera_client.get(route:="/devices"+l.device.id)
  return Light{
    dirigera_client: l.dirigera_client,

  }
}
