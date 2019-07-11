
#!/bin/bash

docker rm -f $(docker ps -aq)
images=( cli orderer.example.com peer0.manufacturerorg.example.com peer1.manufacturerorg.example.com peer0.retailorg.example.com peer1.retailorg.example.com peer0.distributororg.example.com peer1.distributororg.example.com peer0.shiporg.example.com peer1.shiporg.example.com )
for i in "${images[@]}"
do
	echo Removing image : $i
  docker rmi -f $i
done

#docker rmi -f $(docker images | grep none)
images=( cli orderer.example.com peer0.manufacturerorg.example.com peer1.manufacturerorg.example.com peer0.retailorg.example.com peer1.retailorg.example.com peer0.distributororg.example.com peer1.distributororg.example.com peer0.shiporg.example.com peer1.shiporg.example.com )
for i in "${images[@]}"
do
	echo Removing image : $i
  docker rmi -f $(docker images | grep $i )
done


docker network prune
docker volume prune
docker images rm