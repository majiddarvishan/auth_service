set -e

cd ..

version=`cat version`
echo "version: $version"

#cp configs/config.json .

docker build  --network=host -t auth_proxy:$version -f deployments/Dockerfile .
