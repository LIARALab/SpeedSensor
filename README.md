# Real-time Gait Speed Evaluation at Home

This repository contains every material to reproduce our kit, described in our article [1].

As mentionned in the paper, you are able to found here: 
* The full code of the device
* A full description of the hardware
* Schema to reproduce the hardware-part (including the 3D-model)
* Data collected during our experiment

# Full code of the device
The **GoLang** program is available in the folder ```code/```.
However, you may have to change some constants depending on your settings (See file ```code/constants.go``` for more informations).

<u>_It is highly recommended to recompile your own executable._</u>

# Description of the hardware

In this section, the full equipment will be explored.

### Electronics
The module is composed of 3 Infrared Proximity Sensors (*IRPS*), sending their data to an *ADS1x15*,
to convert the analog values to digital ones (which can be read by the Raspberry Pi Zero W).
The protocol used between the *ADS1x15* and the Raspberry Pi is I²C.
The schema is available below and in the folder ```hardware/```.

![Hardware Schematics][hardware_schema]

Here is the devices required to build this module. We use the maker's link, to let you choose your own provider.
* 1x Raspberry Pi Zero W (<https://www.raspberrypi.org/products/raspberry-pi-zero-w/>)
* 1x ADS1x15 (<https://www.adafruit.com/product/1085>)
* 3x Infrared Proximity Sensors (<https://www.sparkfun.com/products/8958>)

### Aesthetics

The 3D-model is provided in the directory ```hardware/```.
Every unit has a _base_ and a _cover_ file.
The main unit is in the file ```hardware/main_unit_(base|cover).stl``` and both surrounding modules are available in ```hardware/left_unit_(base|cover).stl``` & ```hardware/right_unit_(base|cover).stl```.
Every case has holes in its side, to be able to connect everything together. To hide wires, we chose to use an Expandable Sleeving Cable Wire (coming through those holes).

The whole module fully developed is shown in the next figure.

![Full device][full_device]

# Dataset

As explained in our paper [1], we also provide our collected data in the following repository :
<https://github.com/LIARALab/Datasets/tree/master/SpeedSensorDataset>

# References

[1] Chapron, Kévin; Bouchard ,Kévin; Gaboury, Sébastien. (2019, September). "Real-time Gait Speed Evaluation at Home". GOODTECHS 2019 - 5th EAI International Conference on Smart Objects and Technologies for Social Good.


[hardware_schema]: https://github.com/LIARALab/SpeedSensor/raw/master/hardware/Schema_device_bw.png "Hardware schematics"
[full_device]: https://github.com/LIARALab/SpeedSensor/raw/master/hardware/Photo_device_bw.png "Hardware Aesthetics"