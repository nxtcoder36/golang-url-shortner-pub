Long_Url: https://domain/long-url

Short_Url: https://domain/short-url

## Environment Variables
PORT=3000

MONGO_URL=mongodb://mohit:mohitpass@localhost:27017/

MONGO_DB_NAME=url-shortner

MONGO_COLLECTION_NAME=urls

DOMAIN_NAME=http://localhost:3000

REDIS_ADDR=localhost:6379

REDIS_TTL=60

REDIS_PASS=mohitpass

#### REDIS_TTL is time to live in seconds
REDIS_TTL=60

## Installing mongodb driver
go get go.mongodb.org/mongo-driver/mongo

go get go.mongodb.org/mongo-driver/mongo/options

## Installing uuid
github.com/google/uuid

## Installing redis
go get github.com/go-redis/redis

## Running database locally
docker-compose up -d

## Connecting to database
mongosh mongodb://mohit:mohitpass@localhost:27017
### if not installed mongosh
#### for arch
yay -S mongosh-bin

## Running database locally with redis
docker-compose -f docker-compose-redis.yaml up -d

## Connecting to database with redis
redis-cli -h localhost -p 6379 -a mohitpass
### if not installed redis-cli
#### for arch
sudo pacman -S redis

## Running the server
task dev

## Deploying to k9s

### Build and push docker image
task docker-build

task docker-push

### Apply secrets
kubectl apply -f mongo-secret.yaml

kubectl apply -f redis-secret.yaml

kubectl apply -f url-shortner-secret.yaml

### Apply PVCs
kubectl apply -f mongo-pvc.yaml

kubectl apply -f redis-pvc.yaml

### Apply Mongo + Redis deployments 
kubectl apply -f mongo-deployment.yaml

kubectl apply -f redis-deployment.yaml

### Deploy apps
#### Make sure your Kubernetes cluster has an Ingress controller running (like NGINX Ingress Controller).
kubectl get pods -n ingress-nginx
#### OR install NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.1/deploy/static/provider/cloud/deploy.yaml

kubectl apply -f url-shortner-deployment.yaml

### Check everything is running
kubectl get pods

kubectl get svc

### If not able to connect api to cluster from outside (if getting error in above command)
kubectl config use-context <cluster-name>
