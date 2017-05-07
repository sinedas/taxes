# Setup

Application requires mysql databases setuped.
Run script taxes.sql to create database

# Run producer

To run producer execute command:

`go run producer.go`

# Producer endpoints

## Set tax

Examples:

`curl http://localhost:3000/settax/Vilnius/2017/0.2` - set tax 0.2 for year 2017 in Vilnius.

`curl http://localhost:3000/settax/Vilnius/2017-02/0.3` - set tax 0.3 for February of 2017 in Vilnius.

`curl http://localhost:3000/settax/Vilnius/2017-02-15/0.4` - set tax 0.4 for 2017-02-15 in Vilnius.

`curl http://localhost:3000/settax/Vilnius/2017-02-16week/0.5` - set tax 0.5 for week starting 2017-02-16 in Vilnius. Full period 2017-02-16 - 2017-02-22.

* Note: Weeks can overlap each other. If getting dates tax, first week always wins.* 

## Get tax

Example:

`curl http://localhost:3000/tax/Vilnius/2017-02-15` - get tax for date 2017-02-15 in Vilnius.

## Upload data from file

Example:

`curl -i -X POST -H "Content-Type: multipart/form-data" -F "upload=@taxes.csv" http://localhost:3000/upload` - upload taxes.csv with date. Examples files are present.

# Run consumer

Example of running consumer:

`go run consumer.go vilnius 2017-02-16`


