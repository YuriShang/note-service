from fastapi import FastAPI

from app.api.routes import router
from app.schemas.schemas import UserRead
from app.utils.users import fastapi_users, auth_backend


app = FastAPI(title="async-user-service")

# /login, /logout
app.include_router(
    fastapi_users.get_auth_router(auth_backend), 
    tags=["auth"]
)

"""
# /forgot-password
app.include_router(
    fastapi_users.get_reset_password_router(),
    tags=["auth"],
)

# /verify
app.include_router(
    fastapi_users.get_verify_router(UserRead),
    prefix="/auth",
    tags=["auth"],
)
"""
app.include_router(router)