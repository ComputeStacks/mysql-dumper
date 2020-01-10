# MySQL Dump Utility

This is a simple program that will use the ComputeStacks API to find all MySQL hosts in your project, and dump each database into an individual file.

This is designed to be used in conjunction with our [MySQL Backup Image](https://github.com/ComputeStacks/docker/tree/master/utilities/mysql-backup) to perform automated database backups to a remote host.

## Configuration

The following executables must be available on the system:
 
* `mysql`
* `mysqldump`
* `tar`
