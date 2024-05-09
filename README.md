# Prometheus Node Exporter Information

This is a utility which connects to a [Prometheus Node Exporter](https://prometheus.io/docs/guides/node-exporter/) instance and dumps out device and other information.

## Options

By default it'll connect to localhost:9100 but you can pass `-ip` to provide an alternative target address and `-port` to change the target port.

## Output

what the utility does is look in pre-determined parts of the output from Prometheus Node Exporter and get device and other information from them, then de-duplicate and print out.

The output is split into "Disk" , "Filesystem", "Network", "Thermal" and "System Info" categories. These may not all return data, depending on what's available.
