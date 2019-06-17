# Real-time Gait Speed Evaluation at Home

This repository contains every material to reproduce our kit, described in our article [1].

As mentionned in the paper, you are able to found here: 
* The full code of the device
* A full description of the hardware
* Schema to reproduce the hardware-part (including the 3D-model)
* Data collected during our experiment

# Full code of the device
It is available in the folder ```code/```. 
While it is programmed in **GoLang**, it has been compiled to run on many platforms (see ```code/cmd/```).

# Description of the hardware
The module is composed of 3 Infrared Proximity Sensors (*IRPS*), sending their data to an *ADS1x15*, to convert the analog values to digital ones (which can be read by the Raspberry Pi Zero W). The protocol used between the *ADS1x15* and the Raspberry Pi is I²C. The schema is available below and in the folder ```hardware/```.



# References

[1] Chapron, Kévin; Bouchard ,Kévin; Gaboury, Sébastien. (2019, September). "Real-time Gait Speed Evaluation at Home". GOODTECHS 2019 - 5th EAI International Conference on Smart Objects and Technologies for Social Good.
