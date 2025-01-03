curl -X POST http://localhost:8082/orders \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-H "Idempotent-Key: unique-key-12345" \
-H "Content-Type: application/json" \
-d '{
  "user_id": 1,
  "product_requests": [
    {
      "product_id": 101,
      "quantity": 2,
      "mark_up": 10.0,
      "discount": 5.0,
      "final_price": 95.0
    },
    {
      "product_id": 102,
      "quantity": 1,
      "mark_up": 15.0,
      "discount": 3.0,
      "final_price": 112.0
    }
  ],
  "quantity": 3,
  "total": 207.0,
  "total_mark_up": 25.0,
  "total_discount": 8.0,
  "status": "created"
}'
