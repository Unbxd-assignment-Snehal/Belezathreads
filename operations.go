package main

const CREATE_PRODUCT_TABLE = `
CREATE TABLE IF NOT EXISTS PRODUCT (
   productID VARCHAR(255) PRIMARY KEY NOT NULL,
   title VARCHAR(255),
   price DECIMAL(10,2),
   description TEXT,
   categoryID INT,
   FOREIGN KEY (categoryID) REFERENCES CATEGORY(categoryID)
);`

const CREATE_IMAGE_TABLE = `
CREATE TABLE IF NOT EXISTS IMAGE (
   imageID SERIAL PRIMARY KEY NOT NULL,
   imagePath VARCHAR(255),
   productID  VARCHAR(255),
   FOREIGN KEY (productID) REFERENCES PRODUCT(productID)
);`

const CREATE_CATEGORY_TABLE = `
CREATE TABLE IF NOT EXISTS CATEGORY (
   categoryID SERIAL PRIMARY KEY NOT NULL,
   category VARCHAR(255),
   parentCategory VARCHAR(255)
);`

const SELECT_ALL_PRODUCTS = `
SELECT P.productID, P.title, P.price, P.description, P.categoryID, I.imagePath
FROM PRODUCT P JOIN IMAGE I ON P.productID = I.productID;`

const GET_PRODUCT = `

SELECT P.productID, P.title, P.price, P.description, P.categoryID, I.imagePath
FROM PRODUCT P JOIN IMAGE I ON P.productID = I.productID WHERE P.productID = I.productID;`
