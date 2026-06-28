# 🔓 Temporary Admin Setup - Instructions

## What Was Done

I've temporarily unlocked the **Create Admin User** endpoint to allow you to create the first admin without needing to already be logged in as an admin.

### Changes Made:

1. **Added Public Endpoint**: `/api/v1/public/admin/setup`
   - Location: `backend/cmd/server/router.go` (line 52)
   - This endpoint bypasses authentication
   - Uses the same handler as the protected admin creation endpoint

2. **Updated Postman Collection**: 
   - Added `[TEMP] Setup First Admin (Public)` endpoint in the Authentication section
   - Includes example request body
   - No authentication required

## 📋 How to Create Your First Admin

### Option 1: Using Postman (Recommended)

1. Open Postman and import the updated `postman_collection.json`
2. Navigate to: **02 - Authentication** → **[TEMP] Setup First Admin (Public)**
3. Modify the request body if needed:
   ```json
   {
       "username": "admin",
       "email": "admin@example.com",
       "password": "securepass123",
       "full_name": "System Administrator",
       "role_id": ""
   }
   ```
4. Click **Send**
5. You should get a `201 Created` response

### Option 2: Using cURL

```bash
curl -X POST http://localhost:8080/api/v1/public/admin/setup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "securepass123",
    "full_name": "System Administrator",
    "role_id": ""
  }'
```

## ✅ After Creating the First Admin

1. **Login as Admin**:
   - Use the **Admin Login** endpoint in Postman
   - Username: `admin` (or whatever you set)
   - Password: `securepass123` (or whatever you set)
   - This will set the `access_token` automatically

2. **Create Additional Admins** (if needed):
   - Now use the protected endpoint: **10 - Admin Management** → **Admin Users Management** → **Create Admin User**
   - This one requires admin authentication

3. **⚠️ IMPORTANT - Remove the Temporary Endpoint**:
   
   For security, you should remove the public setup endpoint after creating your first admin.

   **To remove it:**
   
   Open `backend/cmd/server/router.go` and delete these lines (52-53):
   ```go
   // TEMPORARY: Public endpoint to create first admin (REMOVE AFTER CREATING FIRST ADMIN)
   e.POST("/api/v1/public/admin/setup", adminHandler.CreateAdminUser)
   ```
   
   Also remove the `adminHandler` parameter from `setupAuthRoutes`:
   ```go
   // Change from:
   func setupAuthRoutes(e *echo.Echo, authHandler *handler.AuthHandler, adminHandler *handler.AdminHandler) {
   
   // Back to:
   func setupAuthRoutes(e *echo.Echo, authHandler *handler.AuthHandler) {
   ```
   
   And update the function call (line 38):
   ```go
   // Change from:
   setupAuthRoutes(e, authHandler, adminHandler)
   
   // Back to:
   setupAuthRoutes(e, authHandler)
   ```

## 🔒 Security Note

The temporary endpoint is **intentionally public** to solve the chicken-and-egg problem of needing an admin to create an admin. However, **it should NOT remain in production code** as it would allow anyone to create admin accounts.

**Remove it as soon as you've created your first admin user!**

## 🧪 Testing

After creating the admin:

1. Test admin login: `POST /api/v1/auth/admin/login`
2. Test protected admin endpoint: `GET /api/v1/admin/me`
3. Verify you can create other admins through the protected route

---

**Need help?** The temporary endpoint uses the exact same validation and creation logic as the protected one, so all your admin creation rules (unique username, unique email, password min length, etc.) still apply.
