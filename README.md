# `netoffload`

`netoffload` is a simple toolkit for simplifying the process of managing network
offloading features on Linux.  It aims at abstracting all of the complexity
involved in managing these features, and providing a simple interface for
enabling and disabling them.

## Usage

Before you can use `offloadctl`, you must first ensure that a few dependencies
are setup on your system.  However, don't worry, `offloadctl` will tell you
exactly what you need to do if it notices that something is missing.

### BIOS configuration

Depending on your hardware, you may need to enable a few options in the BIOS
in order to get SR-IOV working.  You should consult with your hardware vendor
to ensure that you have all the options but you can use the following as a
guideline.

#### Supermicro

You need to make sure that you have the latest BIOS installed on your system
and that you have the following options enabled:

* Advanced > ACPI settings > PCI AER Support > Enabled (if you're using a
  H12/H11 series motherboard with Rome CPU such as EPYC 7xx2)
* Advanced > CPU configuration > SVM Mode > Enabled
* Advanced > NB Configuration > IOMMU > Enabled
* Advanced > NB Configuration > ACS Enable > Enabled
* Advanced > PCIe/PCI/PnP Configuration > SR-IOV Support > Enabled

### Kernel configuration

You must have the following kernel options enabled in your command line:

- `iommu=pt`
- `intel_iommu=on` (for systems with Intel CPUs)
- `amd_iommu=on` (for systems with AMD CPUs)

The most common way to do this is to edit `/etc/default/grub` and add them to
the `GRUB_CMDLINE_LINUX_DEFAULT` variable, then run `update-grub` to update your
grub configuration.

### NIC configuration

There are a few steps necessary to configure your NIC for hardware acceleration.
The exact steps will vary depending on your NIC vendor and model, but the
following can be used as a guideline.

#### Mellanox

1. Install the `mstflint` tools on the compute node which will be used:

   ```bash
   sudo apt-get install mstflint
   ```

1. Get the device's PCI by using `lspci`.

   ```console
   $ lspci | grep Mellanox
   61:00.0 Ethernet controller: Mellanox Technologies MT2892 Family [ConnectX-6 Dx]
   61:00.1 Ethernet controller: Mellanox Technologies MT2892 Family [ConnectX-6 Dx]
   ```

1. Check if SR-IOV is enabled in the firmware

   ```console
   $ sudo mstconfig -d 61:00.0 q | grep SRIOV_EN
        SRIOV_EN                            True(1)
   ```

   If SR-IOV is not enabled, you can enable it with the following command:

   ```console
   $ sudo mstconfig -d 61:00.0 set SRIOV_EN=1
   ```

1. Configure the needed number of VFs

   ```console
   $ sudo mstconfig -d 61:00.0 set NUM_OF_VFS=16
   ```

1. Restart the system

   > **Note**
   >
   > A useful tip is to prefix this command with an extra space (before `sudo`),
   > so that it is not saved in the shell history and prevents accidental reboot.

   ```
   $ sudo reboot
   ```
