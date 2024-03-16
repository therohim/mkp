# Author
### Nama : Nur Rohim
### Email : nrohim893@gmail.com

## Cloning the repository
To clone the repository run the following command:
```bash
$ git clone https://github.com/therohim/mkp.git
```

## Migrate the database
run command below in your terminal. this command will be running the sql file in migrate folder
```bash
migrate -database "postgresql://youruser:yourpass@localhost:5432/yourdatabase?sslmode=disable" -path "./migrate/" up
```
 
## How To Run the Project
copas .env below
 ```bash
APP_VERSION="0.0.1-dev"
APP_SECRET=3b44b7072a761df9769df87cf598e87f

POSTGRE_HOST=localhost
POSTGRE_PORT=5432
POSTGRE_USER=yourusername
POSTGRE_PASS=yourpass
POSTGRE_DB=yourdb
POSTGRE_SSL=disable
POSTGRE_POOL_MIN=10
POSTGRE_POOL_MAX=100
POSTGRE_MAX_LIFE_TIME=60
POSTGRE_MAX_IDLE=10
 ```
 Run in your terminal this command
 ```bash
 make server-start
 ```
 or 
 ```bash
 go run main.go
 ```

 ### documentation postman
 ```bash
 https://documenter.getpostman.com/view/3483861/2sA2xnx9od
 ```

 ### ERD
 ![plot](./assets/image1.png)

 ### Flowchart
 ![plot](./assets/image2.png)