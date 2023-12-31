# Chen's Bond: MongoDB Sharded Cluster Analysis Tool

**Chen's Bond** is a MongoDB analysis tool designed for the purpose of database management and optimization. It primarily operates by querying the 'config' database, providing users with valuable insights into their MongoDB deployment. The tool is particularly focused on improving the performance and efficiency of MongoDB sharded clusters.  To assess a cluster, connect Bond to a mongos instance. Alternatively, it can establish a direct connection to the config database or be used in a cluster restore from a config database backup.

## Build Instructions
Build instructions and user's guide are available at [![Chen's Bond - MongoDB Sharded Cluster Analysis Tool](https://img.youtube.com/vi/equz1z0igv0/0.jpg)](https://youtu.be/equz1z0igv0).

## Changes
### v0.2.0
- Added *Chunk Move Errors*

![Chunk Move Errors](docs/chunk_move_errors.png)

## License
[Apache-2.0 License](LICENSE)
