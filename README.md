# MySQL Dump Utility

This is a simple program that will use the ComputeStacks API to find all MySQL hosts in your project, and dump each database into an individual file.

This is designed to be used in conjunction with our [MySQL Backup Image](https://github.com/ComputeStacks/docker/tree/master/utilities/mysql-backup) to perform automated database backups to a remote host.

## Configuration

The following environmental variables must be defined in order for this to work:

  * `API_HOST`: The ComputeStacks API Endpoint _(e.g. https://portal.example.com/api)_
  * `API_KEY`: Your ComputeStacks API Key
  * `API_SECRET`: Your ComputeStacks API Secret
  * `PROJECT_ID`: Your ComputeStacks Project ID
  
 Additionally, the following executables must be available on the system:
 
   * `mysql`
   * `mysqldump`
   * `tar`
   
   
 