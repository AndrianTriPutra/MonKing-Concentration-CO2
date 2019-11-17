# MonKing-Concentration-CO2

[Device]\
Device Prepare:\
  SBC,gps ublox serie gps+glonas, MHZ19B sensor\
-check port connection with dmesg and change port on the code\
-go to dir /device and run with go run *.go\

[SavetoDB]\
-make DB in mysql with name 'db-monco2'\
-add table with name 'gondril'\
-add colom \
Tanggal	varchar(255) PK\
Jam	varchar(255) PK\
Heading	varchar(255) PK\
Latitude	varchar(255) PK\
Longitude	varchar(255) PK\
Speed	varchar(255) PK\
Temperature	varchar(255) PK\
CO2	varchar(255) PK\

*note\
maybe you can manage that data better than me  \

[webservice]\
-I don't know how to query DB for sortir data from update to old\
-I try by date, but not effective [*maybe i have to add colum number for sortir]\
-I don't know how to send many data  effectively from go to html\

I hope could help you\
and\
I think you can you better\

Thanks
