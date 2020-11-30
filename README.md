# Dafflabs

Experimental service for fun and profit.

## Todo
1. Create table to store product state
    - url
    - created_at
    - updated_at
    - original_price
    - current_price
2. Create wishlist table
    - product_url
    - created_at
    - customer_ref_id
3. Set scheduler to check product price regularly 
    - get price from source
    - get last fetched price from table
    - if new price < last fetched price
        - update table
        - notify to subscribers
4. Create API to create wishlist
5. Create API to retrieve wishlist
6. Create webhook in seanmcapp

## Roadmap

1. Get price alert from Indonesian e-commerce and send it via seanmcapp
2. Get news page urls and send it via seanmcapp
