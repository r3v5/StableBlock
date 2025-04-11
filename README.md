# **StableBlock**
★ StableBlock is REST API written in Go, deployed using k8s and simulates lightweight blockchain operations and protected with JWT

<a name="readme-top"></a>

  <h3 align="center">StableBlock API</h3>
  </p>
</div>






<!-- ABOUT THE PROJECT -->
## About The Project

StableBlock is REST API written in Go, deployed to local Kubernetes cluster that runs behind Docker Desktop and simulates lightweight blockchain operations

## StableBlock Database Overview
![System Design](https://raw.githubusercontent.com/r3v5/StableBlock/main/StableBlock-DB-architecture.png)



## Deployment & Infrastructure
•  **Docker Image based on golang:1.24.2-alpine3.21**

•  **Local Kubernetes cluster**: k8s cluster that runs behind Docker Desktop

•  **postgres service**: ClusterIP with 5432 port

•  **stableblock-api service**: NodePort 30080

•  **postgres deployment with 1 replica**

•  **stableblock-api deployment with 10 replicas**

•  **stableblock namespace**

•  **Two Makefiles**: Makefile.k8s and Makefile.local

### Built With

 <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=go,k8s,docker,postgres,linux" />
  </a>

<p align="right">(<a href="#about-the-project">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Installation and deployment to k8s

1. Clone the repo ✅
  ```sh
   https://github.com/r3v5/StableBlock.git
   ```
2. Navigate to the project directory ✅
  ```sh
   cd StableBlock
   ```
3. Create a .env file local development ✅
  ```
	DB_HOST=localhost
	DB_PORT=5432
	DB_USER=postgres
	DB_PASSWORD=postgres
	DB_NAME=postgres
	DB_SSLMODE=disable
	DB_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
	JWT_SECRET=buterinthegodnow
```
	
4. Build docker image ✅
  ```
   docker build -f infra/docker/Dockerfile -t stableblock-api:v1 .                                                                                                         
  ``` 
 5. Verify image is created ✅
 ```
 docker images     
 ```  
  4. In folder infra/k8s create configmap.yml (config file for k8s deployment) ✅
   ```
apiVersion: v1
kind: ConfigMap
metadata:
	name: stableblock-config
	namespace: stableblock
data:
     DB_HOST: "postgres"
     DB_PORT: "5432"
	 DB_NAME: "postgres"
	 DB_SSLMODE: "disable"
   ```
  
  5. In folder infra/k8s create secret.yml (config file for k8s deployment) ✅
  ```
apiVersion: v1
kind: Secret
metadata:
	name: stableblock-secrets
	namespace: stableblock
type: Opaque
stringData:
	DB_USER: "postgres"
	DB_PASSWORD: "postgres"
	JWT_SECRET: "buterinthegodnow"
	DB_DSN: "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
   ```
 6. Create stableblock namespace in Kubernetes cluster ✅
  ```
kubectl create -f infra/k8s/namespace.yml                                                                                                                                   

namespace/stableblock created
   ```

 7. Create PersistentVolumeClaim for db in Kubernetes cluster ✅
  ```
kubectl create -f infra/k8s/postgres-pvc.yml   
                                                                                                                             
persistentvolumeclaim/postgres-pvc created                                                                                                                                 
   ```
 8. Create ConfigMap in Kubernetes cluster ✅
  ```
kubectl create -f infra/k8s/configmap.yml  
                                                                                                                                 
configmap/stableblock-config created                                                                                                                             
   ```   
   9. Create Secret in Kubernetes cluster ✅
  ```
kubectl create -f infra/k8s/secret.yml      
                                                                                                                                
secret/stableblock-secrets created                                                                                                                           
   ```  
  10. Create Postgres Service in Kubernetes cluster ✅
  ```
kubectl create -f infra/k8s/postgres-service.yml                                                                                                                       

service/postgres created                                                                                                                         
   ```  
   11. Make deployment of Postgres Service in Kubernetes cluster ✅
  ```
 kubectl create -f infra/k8s/postgres-deployment.yml                              
                                                                                            
deployment.apps/postgres created                                                                                                                       
   ```  
  12. Verify that stuff with Postgres was created in Kubernetes cluster ✅
  ```
 kubectl get all -n stableblock                                                                                                                                              

NAME                            READY   STATUS    RESTARTS   AGE
pod/postgres-6d9d5894dc-pdzp6   1/1     Running   0          53s

NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/postgres   ClusterIP   10.110.210.175   <none>        5432/TCP   6m41s

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/postgres   1/1     1            1           53s

NAME                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/postgres-6d9d5894dc   1         1         1       53s                                                                                                                 
   ```  
13. Create NodePort Service in Kubernetes cluster ✅
  ```
  kubectl create -f infra/k8s/stable-block-api-service.yml
                                                                                                                      
service/stableblock-api created                                                                                                                     
   ``` 
   14. Finally make a deployment of NodePort service with 10 replicas in Kubernetes cluster ✅
  ```
 kubectl create -f infra/k8s/stable-block-api-deployment.yml     
                                                                                                             
deployment.apps/stableblock-api created                                                                                                                   
   ``` 
   15. Get all pods, services and deployments for given namespace Kubernetes cluster ✅
  ```
 kubectl get all -n stableblock                                                                                                                                              

NAME                                   READY   STATUS    RESTARTS   AGE
pod/postgres-6d9d5894dc-pdzp6          1/1     Running   0          13m
pod/stableblock-api-6bf4465676-4x24d   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-5mgvq   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-9fxf7   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-bx8rn   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-cskdp   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-cwscw   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-h6scp   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-ng4qn   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-t4zvv   1/1     Running   0          8m34s
pod/stableblock-api-6bf4465676-td6wc   1/1     Running   0          8m34s

NAME                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/postgres          ClusterIP   10.110.210.175   <none>        5432/TCP         18m
service/stableblock-api   NodePort    10.110.46.133    <none>        8080:30080/TCP   9m53s

NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/postgres          1/1     1            1           13m
deployment.apps/stableblock-api   10/10   10           10          8m34s

NAME                                         DESIRED   CURRENT   READY   AGE
replicaset.apps/postgres-6d9d5894dc          1         1         1       13m
replicaset.apps/stableblock-api-6bf4465676   10        10        10      8m34s
                                                                                                                 
 ``` 

16. List the pods, select any from stable-block-api and run migrations on it ✅
  ```
kubectl exec -it stableblock-api-6bf4465676-4x24d -n stableblock -- make -f Makefile.k8s migrate-up                                                                         

migrate -path database/migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" up
20250406135230/u init (15.89675ms)
20250411124832/u add_date_created_to_accounts (19.840417ms)                                                                                                                   
   ``` 
  17. Then verify the current migration version ✅
  ```
kubectl exec -it stableblock-api-6bf4465676-4x24d -n stableblock -- make -f Makefile.k8s migrate-version         
                                                          
migrate -path database/migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" version
20250411124832                                                                                                                 
   ``` 
  18. All is good, then call welcome endpoint in Postman or using curl and explore docs to see all the endpoints ✅
  ```
curl --location 'http://localhost:30080/api/v1/welcome/'

{
	"message":  "Welcome to StableBlock 👋"
}                                                                                                                
   ``` 
   
### Postman collection with documentation for API
Published Postman collection: [API docs](https://documenter.getpostman.com/view/27242366/2sB2cYdLba)

<!-- CONTACT -->
## Contact

Ian Miller - [linkedin](https://www.linkedin.com/in/ian-miller-620a63245/) 

Project Link: [https://github.com/r3v5/StableBlock](https://github.com/r3v5/StableBlock)

<p align="right">(<a href="#about-the-project">back to top</a>)</p>