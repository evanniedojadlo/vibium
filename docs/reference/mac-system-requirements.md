# System Requirements

Hardware requirements for running two macOS guest VMs simultaneously with Xcode, iOS Simulator, Android Studio, and browser testing.

## Host Machine (Apple Silicon Mac)

|               | Minimum | Recommended |
| ------------- | ------- | ----------- |
| **RAM**       | 32 GB   | 64 GB       |
| **SSD**       | 512 GB  | 1 TB        |

## Each macOS Guest VM

|                    | Minimum | Recommended |
| ------------------ | ------- | ----------- |
| **RAM**            | 10 GB   | 16 GB       |
| **Virtual Disk**   | 192 GB  | 256 GB      |

## Assumed Guest Workload

- macOS (base install)
- Xcode + iOS Simulator runtimes
- Android Studio + SDK + emulator system images
- Chrome and Firefox
- Development workspace and build caches

## Assumed Host Workload

- Two macOS guest VMs running simultaneously
- Xcode and iOS Simulator on the host

## Notes

- Below 32 GB host RAM, the two-VM workflow is not supported. Drop to one VM or run without VM isolation.
- 24 GB host RAM can run a single VM (10 GB) with Xcode and iOS Simulator on the host, but this is not a recommended configuration.
- Guest virtual disks are sparse — a 256 GB virtual disk only consumes actual space as it fills, so oversizing has minimal cost on hosts with sufficient storage.
- After provisioning, monitor Activity Monitor's Memory Pressure graph (green = fine, yellow = tight, red = undersized) and Swap Used to verify the configuration is adequate.
