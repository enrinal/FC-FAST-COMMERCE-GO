import http from 'k6/http';
import { check, sleep } from 'k6';

// Define the URL of the Order Service
const ORDER_SERVICE_URL = 'http://localhost:8082/orders';

// Define the JWT token and Idempotent Key (replace with actual values)
const JWT_TOKEN = 'YOUR_JWT_TOKEN';

// Create random IDEMPOTENT_KEY
const randomString = () => Math.random().toString(36).substring(2);
const IDEMPOTENT_KEY = randomString();

// Function to generate a random order
export default function () {
  const payload = JSON.stringify({
    user_id: 1,
    product_requests: [
      {
        product_id: 101,
        quantity: Math.floor(Math.random() * 5) + 1,  // Random quantity between 1-5
        mark_up: 10.0,
        discount: 5.0,
        final_price: 95.0
      },
      {
        product_id: 102,
        quantity: Math.floor(Math.random() * 5) + 1,  // Random quantity between 1-5
        mark_up: 15.0,
        discount: 3.0,
        final_price: 112.0
      }
    ],
    quantity: 3,
    total: 207.0,
    total_mark_up: 25.0,
    total_discount: 8.0,
    status: 'created'
  });

  const params = {
    headers: {
      'Authorization': `Bearer ${JWT_TOKEN}`,
      'Idempotent-Key': IDEMPOTENT_KEY,
      'Content-Type': 'application/json',
    },
  };

  // Send the POST request to the Order Service
  let response = http.post(ORDER_SERVICE_URL, payload, params);

  // Check if the status is 200 (OK)
  check(response, {
    'is status 200': (r) => r.status === 200,
  });

  // Simulate some thinking time between requests
  sleep(1);
}
