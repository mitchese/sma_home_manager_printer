# SMA Home Manager printer

This simple program prints the values of the [SMA Home Manager](https://www.sma.de/en/products/monitoring-control/sunny-home-manager-20.html)

I wrote this as an initial test before implementing a fake ET340 in a Victron system. Use this at your own risk, I have
no association with Victron but am providing this for anyone who already has an SMA Home Manager (and no room left in their fuse box for a further meter)

## How it works

The SMA Home Manager publishes every second a UDP message which includes all the values, both instantaneous usage as well as counters.

In order to 'fake' an ET340 in the Victron system, this program will subscribe to the updates of the meter and inject these into dbus.

This program subscribes and decodes the important values, which are used to test if the connection works. It prints these out for manual verification

To run, just run in a terminal

## Testing on Victron Venus OS

You can use this program to verify that your Venus OS can connect to an SMA meter, instead of using the ET340/EM24 meters.

I have tested this on the Victron Venus GX, but the same should work for any of the Victron GX devices. In the menu:

* Activate "Superuser" under Settings -> General (see https://www.victronenergy.com/live/ccgx:root_access)
* Set an SSH Password and enable SSH on LAN
* Download the `sma_home_manager_printer.armv7l` to your computer
* Use WinSCP or other scp program to copy `sma_home_manager_printer.armv7l` to your Venus
* Use Putty (or similar) to ssh to the Venus GX. Run with `./sma_home_manager_printer.armv7l`

Output should resemble the following:
```
2020/05/27 20:47:45 -----------------------------------------------------
2020/05/27 20:47:45 Received datagram from meter
2020/05/27 20:47:45 Uid:  253910
2020/05/27 20:47:45 Serial:  3004906401
2020/05/27 20:47:45 Total W:  484
2020/05/27 20:47:45 Total Buy kWh:   6675.477900000001
2020/05/27 20:47:45 Total Sell kWh:  3200.4507000000003
2020/05/27 20:47:46 +-----+-----------+-----------+-----------+
2020/05/27 20:47:45 |value|   L1      |     L2    |   L3      |
2020/05/27 20:47:45 +-----+-----------+-----------+-----------+
2020/05/27 20:47:45 |  V  |   231.18  |   230.75  |   231.35  |
2020/05/27 20:47:45 |  A  |    -0.11  |     0.52  |     1.68  |
2020/05/27 20:47:45 |  W  |   -24.90  |   120.90  |   388.10  |
2020/05/27 20:47:45 | kWh |  3084.11  |  1213.50  |  3002.64  |
2020/05/27 20:47:45 | kWh |  1309.15  |  1599.17  |   916.91  |
2020/05/27 20:47:45 +-----+-----------+-----------+-----------+
2020/05/27 20:47:46 -------------------------------------------
2020/05/27 20:47:46 Received datagram from meter
2020/05/27 20:47:46 Uid:  253910
2020/05/27 20:47:46 Serial:  3004906401
2020/05/27 20:47:46 Total W:  490.6
2020/05/27 20:47:46 Total Buy kWh:   6675.478
2020/05/27 20:47:46 Total Sell kWh:  3200.4507000000003
2020/05/27 20:47:46 +-----+-----------+-----------+-----------+
2020/05/27 20:47:46 |value|   L1      |     L2    |   L3      |
2020/05/27 20:47:46 +-----+-----------+-----------+-----------+
2020/05/27 20:47:46 |  V  |   231.57  |   230.76  |   231.03  |
2020/05/27 20:47:46 |  A  |    -0.11  |     0.55  |     1.68  |
2020/05/27 20:47:46 |  W  |   -24.80  |   127.60  |   387.80  |
2020/05/27 20:47:46 | kWh |  3084.11  |  1213.50  |  3002.64  |
2020/05/27 20:47:46 | kWh |  1309.15  |  1599.17  |   916.91  |
2020/05/27 20:47:46 +-----+-----------+-----------+-----------+
```

You can verify the Serial number with the number shown on your SMA inverter under Device Configuration. The rest of
the values should make sense.

## Problems

If this program begins outputting data as above, but stops within a few minutes, ensure that "IGMP Snooping" is enabled
on your network switches/routers (whatever's between this and the home manager). This is necessary for multicast.
I did not have this enabled on my switch, _somehow_ the SMA inverters did not have an issue (they must have constantly
re-connected). This program should run forever until stopped.

