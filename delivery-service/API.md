# Food Delivery Service API Documentation

- Global Error Response
    ```json
    {
        "error": "abcdxyz"
    }
    ```

## User Service
### Authentication
- `POST /auth/register`: Register new user
    - Request
    ```json
    {
        "email": "duong@example.com",
        "password": "supersecurepassword",
        "name": "Thai Duong",
        "gender": "gay",
        "phone": "1234567890",
        "role": "customer" // "customer", "shipper"
    }
    ```
    - Response
    ```json
    {
        "userId": "dsfsdfjk",
        "email": "duong@example.com",
        "password": "supersecurepassword",
        "name": "Thai Duong",
        "gender": "gay",
        "phone": "1234567890",
        "role": "customer" // "customer", "shipper"
    }
    ```

- `POST /auth/login`: Login
    - Request
    ```json
    {
        "email": "user@example.com",
        "password": "securepassword"
    }
    ```
    - Response
    ```json
    {
        "userId": 01,
        "accessToken": "asfjhasdkjfhsk",
        "role": "customer"
    }
    ```

### User Profile Management
- `GET users/profile`: Get current user profile
    - Headers
        - Authorization: Bearer (token) (all roles)
    - Response:
    ```json
    {
        "userId": 01,
        "email": "duong@example.com",
        "name": "Thai Duong",
        "gender": "gay",
        "phone": "1234567890",
        "role": "customer" // "customer", "shipper"
    }
    ```
- `PATCH users/profile`: Update current user profile
    - Headers
        - Authorization: Bearer (token) (all roles)
    - Request:
    ```json
    {
        "email": "duong@example.com",
        "name": "Thai Duong",
        "gender": "lesbian",
        "phone": "1234567890"
    }
    ```
    - Response:
    ```json
    {
        "email": "duong@example.com",
        "name": "Thai Duong",
        "gender": "lesbian",
        "phone": "1234567890"
    }
    ```
- `PUT users/password`: Update user password
    - Headers
        - Authorization: Bearer (token) (all roles)
    - Request
    ```json
    {
        "currentPassword": "oldpassword",
        "newPassword": "newpassword"
    }
    ```
    - Response
    ```json
    {
        "message": "password changed successfully"
    }
    ```

## Restaurant Service
### Restaurant Management
- `GET /restaurant`: Get restaurant info
    - Response
    ```json
    {
        "name": "Coastal Breeze Bistro",
        "description": "A charming seaside restaurant offering fresh seafood and Mediterranean-inspired dishes in a relaxed atmosphere with ocean views.",
        "address": "Chung cư Season Avenue",
        "phoneNumber": "+1-415-555-8723",
        "isActive": true,
        "openTime": "0000-01-01T11:00:00Z",
        "closeTime": "0000-01-01T23:00:00Z"
    }
    ```
- `PATCH /restaurant`: Update restaurant details
    - Headers
        - Authorization: Bearer (token) (`admin` role)
    - Request
        ```json
        {
            "name": "Coastal Breeze Bistro",
            "description": "A charming seaside restaurant offering fresh seafood and Mediterranean-inspired dishes in a relaxed atmosphere with ocean views.",
            "address": "Chung cư Season Avenue",
            "phoneNumber": "+1-415-555-8723",
            "isActive": true,
            "openTime": "0000-01-01T11:00:00Z",
            "closeTime": "0000-01-01T23:00:00Z"
        }
        ```
    - Response
        ```json
        {
            "name": "Coastal Breeze Bistro",
            "description": "A charming seaside restaurant offering fresh seafood and Mediterranean-inspired dishes in a relaxed atmosphere with ocean views.",
            "address": "Chung cư Season Avenue",
            "phoneNumber": "+1-415-555-8723",
            "isActive": true,
            "openTime": "0000-01-01T11:00:00Z",
            "closeTime": "0000-01-01T23:00:00Z"
        }
        ```
### Menu Management
- `POST /restaurant/menu/item`: Add new menu items
    - Headers
        - Authorization: Bearer (token) (`admin` role)
    - Request
        ```json
        {
            "name": "Pan-Seared Salmon",
            "description": "Fresh Atlantic salmon seared to perfection, served with roasted vegetables and lemon-dill sauce.",
            "price": 24000,
            "isAvailable": true
        }
        ```
    - Response
        ```json
        {
            "message": "menu retrieved successfully"
        }
        ```
- `GET /restaurant/menu`: Get menu
    - Response
    ```json
    [
        {
            "id": 1,
            "name": "Pan-Seared Salmon",
            "description": "Fresh Atlantic salmon seared to perfection, served with roasted vegetables and lemon-dill sauce.",
            "price": 24000,
            "isAvailable": true
        } 
    ]
    ```
- `PATCH /restaurant/menu/item/{itemId}`: Add new menu items
    - Headers
        - Authorization: Bearer (token) (`admin` role)
    - Request
        ```json
        {
            "name": "Pan-Seared Salmon",
            "description": "Fresh Atlantic salmon seared to perfection, served with roasted vegetables and lemon-dill sauce.",
            "price": 24000,
            "isAvailable": true
        }
        ```
    - Response
        ```json
        {
            "id": 1,
            "name": "Pan-Seared Salmon",
            "description": "Fresh Atlantic salmon seared to perfection, served with roasted vegetables and lemon-dill sauce.",
            "price": 24000,
            "isAvailable": true
        } 
        ```
- `DELETE /restaurant/menu/item/{itemId}`: Delete menu items
    - Headers
        - Authorization: Bearer (token) (`restaurant_owner` role)
    - Response
        ```json
        {
            "message": "Menu item deleted successfully"
        }
        ```

## Delivery Service
### Shipper Management
- `POST /shippers`: Create new shipper
    - Headers
        - Authorization: Bearer (token) (`admin` role)
    - Request
    ```json
    {
        "email": "duong@example.com",
        "password": "supersecurepassword",
        "name": "Thai Duong",
        "gender": "gay",
        "phone": "1234567890",
        "vehicleType": "car",
        "vehiclePlate": "30K-999.99"
    }
    ```
    - Response
    ```json
    {
        "userId": 1,
        "email": "duong@example.com",
        "name": "Thai Duong",
        "gender": "gay",
        "phone": "1234567890",
        "role": "shipper",
        "vehicleType": "car",
        "vehiclePlate": "30K-999.99",
        "totalDeliveries": 0,
        "status": "available"
    }
    ```
- `GET /shippers`: Get list of shipper
    - Headers
        - Authorization: Bearer (token) (`admin` role)
    - Query Parameters
        - limit (optional, default 10): Maximum number of restaurants to return.
        - offset (optional, default 0): Number of restaurants to skip for pagination.
    - Response
    ```json
    [
        {
            "userId": 1,
            "email": "duong@example.com",
            "name": "Thai Duong",
            "gender": "gay",
            "phone": "1234567890",
            "role": "shipper",
            "vehicleType": "car",
            "vehiclePlate": "30K-999.99",
            "totalDeliveries": 0,
            "status": "available"
        }
    ]
    ```
- `GET /shippers/{userId}`: Get shipper information
    - Headers
        - Authorization: Bearer (token) (`shipper` and `admin` role)
    - Response
    ```json
    {
        "userId": 1,
        "email": "duong@example.com",
        "name": "Thai Duong",
        "gender": "gay",
        "phone": "1234567890",
        "role": "shipper",
        "vehicleType": "car",
        "vehiclePlate": "30K-999.99",
        "totalDeliveries": 0,
        "status": "available"
    }
    ```
### Delivery Management
- `POST /delivery`: Create new deliveries (only for internal uses)
    - Headers
        - Authorization: Bearer (token) (`customer` role)
    - Request
    ```json
    {
        "orderId": 1,
        "restaurantAddress": "Season Avenue",
        "shippingAddress": "PTIT"
    }
    ```
    - Response
    ```json
    {
        "deliveryId": 1,
        "orderId": 1,
        "distance": 1.410828,
        "duration": 4.85605,
        "fee": 7054,
        "fromCoords": [
            105.786744,
            20.986808
        ],
        "toCoords": [
            105.78992,
            20.981995
        ],
        "geometryLine": "q~a_CcntdSeAzD{AtEzHrKbA_A?}B\\gEvBwIxBwEtDmE|G}G",
        "status": "assigned",
        "shipper": {
            {
                "userId": 1,
                "email": "duong@example.com",
                "name": "Thai Duong",
                "gender": "gay",
                "phone": "1234567890",
                "role": "shipper",
                "vehicleType": "car",
                "vehiclePlate": "30K-999.99",
                "totalDeliveries": 0,
                "status": "available"
            }
        }
    }
    ```
- `PUT /delivery/{deliveryId}/status`: Update delivery status
    - Headers
        - Authorization: Bearer (token) (`shipper` and `admin` role)
    - Request
    ```json
    {
        "status": "delivering"
    }
    ```
    - Response
    ```json
    {
        "message": "Delivery status updated successfully"
    }
- `GET /delivery/shipper/{shipperId}`: Get delivery based on shipperId
    - Headers
        - Authorization: Bearer (token) (`shipper` and `admin` role)
    - Query Parameters
        - limit
        - offset
    - Response
    ```json
    [
        {
            "deliveryId": 1,
            "orderId": 1,
            "distance": 1.410828,
            "duration": 4.85605,
            "fee": 7054,
            "fromCoords": [
                105.786744,
                20.986808
            ],
            "toCoords": [
                105.78992,
                20.981995
            ],
            "geometryLine": "q~a_CcntdSeAzD{AtEzHrKbA_A?}B\\gEvBwIxBwEtDmE|G}G",
            "status": "assigned"
        }
    ]
    ```
- `GET /delivery/order/{orderId}`: Get delivery info for an order
    - Headers
        - Authorization: Bearer (token) (`customer` and `admin` role)
    - Response
    ```json
    {
        "deliveryId": 1,
        "orderId": 1,
        "distance": 1.410828,
        "duration": 4.85605,
        "fee": 7054,
        "fromCoords": [
            105.786744,
            20.986808
        ],
        "toCoords": [
            105.78992,
            20.981995
        ],
        "geometryLine": "q~a_CcntdSeAzD{AtEzHrKbA_A?}B\\gEvBwIxBwEtDmE|G}G",
        "status": "assigned",
        "shipper": {
            {
                "userId": 1,
                "email": "duong@example.com",
                "name": "Thai Duong",
                "gender": "gay",
                "phone": "1234567890",
                "role": "shipper",
                "vehicleType": "car",
                "vehiclePlate": "30K-999.99",
                "totalDeliveries": 0,
                "status": "available"
            }
        }
    }
    ```
## OrderService
### Order Management
- `POST /orders`: Create new order
    - Headers
        - Authorization: Bearer (token) (`customer` role)
    - Request
    ```json
    {
        "shippingAddress": "22 Ao Sen, Mộ Lao, Hà Đông, Hà Nội",
        "phoneNumber": "1234567890",
        "orderItems": [
            {
                "menuItemId": 1,
                "quantity": 3
            },
            {
                "menuItemId": 2,
                "quantity": 3
            },
            {
                "menuItemId": 3,
                "quantity": 3
            }
        ]
    }
    ```
    - Response
    ```json
    {
        "id": 3,
        "userId": 1,
        "shippingAddress": "22 Ao Sen, Mộ Lao, Hà Đông, Hà Nội",
        "phoneNumber": "1234567890",
        "status": "created",
        "subtotal": 249000,
        "deliveryFee": 7054,
        "totalAmount": 256054,
        "orderItems": [
            {
                "menuItemId": 1,
                "quantity": 3,
                "unitPrice": 24000,
                "totalPrice": 72000
            },
            {
                "menuItemId": 2,
                "quantity": 3,
                "unitPrice": 27000,
                "totalPrice": 81000
            },
            {
                "menuItemId": 3,
                "quantity": 3,
                "unitPrice": 32000,
                "totalPrice": 96000
            }
        ],
        "delivery": {
            {
                "distance": 1.410828,
                "duration": 4.85605,
                "fee": 7054,
                "fromCoords": [
                    105.786744,
                    20.986808
                ],
                "toCoords": [
                    105.78992,
                    20.981995
                ],
                "geometryLine": "q~a_CcntdSeAzD{AtEzHrKbA_A?}B\\gEvBwIxBwEtDmE|G}G",
                "status": "assigned",
                "shipper": {
                    {
                        "name": "Thai Duong",
                        "gender": "gay",
                        "phone": "1234567890",
                        "vehicleType": "car",
                        "vehiclePlate": "30K-999.99",
                    }
                }
            }
        }
    }
    ```
- `GET /orders/{orderId}`: Get order by id
    - Headers
        - Authorization: Bearer (token) (`customer` and `admin` role)
    - Response
    ```json
    {
        "id": 3,
        "userId": 1,
        "shippingAddress": "22 Ao Sen, Mộ Lao, Hà Đông, Hà Nội",
        "phoneNumber": "1234567890",
        "status": "created",
        "subtotal": 249000,
        "deliveryFee": 7054,
        "totalAmount": 256054,
        "orderItems": [
            {
                "menuItemId": 1,
                "quantity": 3,
                "unitPrice": 24000,
                "totalPrice": 72000
            },
            {
                "menuItemId": 2,
                "quantity": 3,
                "unitPrice": 27000,
                "totalPrice": 81000
            },
            {
                "menuItemId": 3,
                "quantity": 3,
                "unitPrice": 32000,
                "totalPrice": 96000
            }
        ],
        "delivery": {
            {
                "distance": 1.410828,
                "duration": 4.85605,
                "fee": 7054,
                "fromCoords": [
                    105.786744,
                    20.986808
                ],
                "toCoords": [
                    105.78992,
                    20.981995
                ],
                "geometryLine": "q~a_CcntdSeAzD{AtEzHrKbA_A?}B\\gEvBwIxBwEtDmE|G}G",
                "status": "assigned",
                "shipper": {
                    {
                        "name": "Thai Duong",
                        "gender": "gay",
                        "phone": "1234567890",
                        "vehicleType": "car",
                        "vehiclePlate": "30K-999.99",
                    }
                }
            }
        }
    }
    ```
- `GET /orders`: Get list of order
    - Headers
        - Authorization: Bearer (token) (`admin` role)
    - Query parameter
        - limit (optional)
        - offset (optional)
        - userId (optional) (can be used with `customer_role`)
    - Response
    ```json
    [
        {
            "id": 3,
            "userId": 1,
            "shippingAddress": "22 Ao Sen, Mộ Lao, Hà Đông, Hà Nội",
            "phoneNumber": "1234567890",
            "status": "created",
            "subtotal": 249000,
            "deliveryFee": 7054,
            "totalAmount": 256054
        }
    ]
    ```
- `PUT /orders/{orderId}/status`: Update user status
    - Headers
        - Authorization: Bearer (token) (all role)
    - Request
    ```json
    {
        "status": "cancelled"
    }
    ```
    - Response
    ```json
    {
        "message": "status updated successfully"
    }
    ```    


