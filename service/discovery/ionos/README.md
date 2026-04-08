# IONOS API

## Depth for Resources 

### Compute – Server Resource  
| depth | Data returned inside the Server JSON |
|-------|--------------------------------------|
| 0 | Only the server URI (`href`) |
| 1 | Server properties + links (`href`, `id`) to child resources (NICs, volumes, CD-ROMs, …) |
| 2 | Server properties **+ full child resources** (complete NIC, Volume, CD-ROM objects) |
| 3 | Same as depth 2 **+ grandchildren resources**, e.g. firewall rules of NICs, snapshots of volumes |

### Storage – Volume Resource  
| depth | Data returned inside the Volume JSON (stand-alone or embedded) |
|-------|---------------------------------------------------------------|
| 0 | Only the volume URI (`href`) |
| 1 | Volume properties + links (`href`, `id`) to snapshots/backups |
| 2 | Volume properties **+ full snapshots & backups** |
| 3 | Volume properties, snapshots & backups **+ their child resources** (e.g. scheduler details) |
