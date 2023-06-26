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

