package devices

import "time"

type Attributes struct {
	custom_name        string
	model              string
	manufacturer       string
	firmware_version   string
	hardware_version   string
	serial_number      string
	ota_status         string
	ota_state          string
	ota_progress       string
	ota_policy         string
	ota_schedule_start time.Time
	ota_schedule_end   time.Time
}

type Capabilities struct {
	can_send    []string
	can_receive []string
}

type Room struct {
	id    string
	name  string
	color string
	icon  string
}

type Device struct {
	id           string
	relation_id  string
	d_type       string
	device_type  string
	created_at   time.Time
	is_reachable bool
	last_seen    time.Time
	attributes   Attributes
	capabilities Capabilities
	room         Room
	device_set   []string
	remote_links []string
	is_hidden    bool
}

func NewDevice(data Device) Device {
	return Device{
		id:           data.id,
		relation_id:  data.relation_id,
		d_type:       data.d_type,
		device_type:  data.device_type,
		created_at:   data.created_at,
		is_reachable: data.is_reachable,
		last_seen:    data.last_seen,
		attributes:   data.attributes,
		capabilities: data.capabilities,
		room:         data.room,
		device_set:   data.device_set,
		remote_links: data.remote_links,
		is_hidden:    data.is_hidden,
	}
}
