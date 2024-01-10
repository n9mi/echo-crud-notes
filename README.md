# Go Echo Notes CRUD
Simple CRUD with golang Echo

## **Packages used**
- github.com/joho/godotenv/cmd/godotenv@latest
- github.com/labstack/echo/v4
- github.com/go-playground/validator/v10
- gorm.io/driver/postgres
- gorm.io/gorm 
- github.com/stretchr/testify

## **Run the migration**
First, create two databases for main usage and testing purpose and configure the .env. When you run the program by
```
go run .
```
Gorm will automatically migrate and seed the database. The testing database will also automatically migrate when you run the testing file.
```
go test -v ./test
```
![image](https://github.com/naomigrain/echo-crud-notes/assets/113373725/6964f429-b16c-4377-9f08-a05d9c6821ad)

## **Structure**
Based on repository pattern, this project use:
- Repository layer: For accessing db in the behalf of project to store/update/delete data
- Service layer: Contains set of logic/action needed to process data/orchestrate those data
- Models layer: Contains set of entity/actual data attribute
- Controller layer: Acts to mapping users input/request and presented it back to user as relevant responses

## **API Endpoints**
You can use <a href="https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi">this VSCode Extension</a> to preview the OpenApi .yml file

