from fastapi import Depends, Request, HTTPException, status, APIRouter
from fastapi_users import exceptions
from fastapi_users.router.common import ErrorCode, ErrorModel

from app.schemas.schemas import UserCreate, UserRead, UsernameUpdate, PasswordUpdate
from app.schemas.responses import password_change_responses, username_change_responses
from app.models.models import User
from app.utils.users import UserManager, current_active_user, get_user_manager


router = APIRouter()


@router.post(
    "/register",
    response_model=UserRead,
    status_code=status.HTTP_201_CREATED,
    tags=["user"]
    )
async def create_user(
    request: Request,
    new_user: UserCreate,
    user_manager: UserManager = Depends(get_user_manager)
) -> UserRead:
    try:
        user = await user_manager.create(new_user, request)
    except exceptions.UserAlreadyExists:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=ErrorCode.REGISTER_USER_ALREADY_EXISTS,
        )
    except exceptions.InvalidPasswordException as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail={
                "code": ErrorCode.REGISTER_INVALID_PASSWORD,
                "reason": e.reason,
            },
        )
    return UserRead.model_validate(user)


@router.get("/me", 
            response_model=UserRead,
            tags=["user"]
)
async def get_user(
    request: Request,
    user: User = Depends(current_active_user),
    user_manager: UserManager = Depends(get_user_manager)
) -> UserRead:
    user = await user_manager.get_by_email(user.username)
    return UserRead.model_validate(user)


@router.patch(
    "/me", 
    response_model=UserRead,
    tags=["user"],
    responses=username_change_responses
)
async def change_username(
    request: Request,
    new_username: UsernameUpdate,
    user: User = Depends(current_active_user),
    user_manager: UserManager = Depends(get_user_manager)
) -> UserRead:
    try:
        user = await user_manager.update(new_username, user, request)
    except exceptions.UserAlreadyExists:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=ErrorCode.REGISTER_USER_ALREADY_EXISTS,
        )
    return UserRead.model_validate(user)


@router.patch(
    "/me/password", 
    response_model=UserRead,
    tags=["user"],
    responses=password_change_responses
)
async def change_password(
    request: Request,
    new_password: PasswordUpdate,
    user: User = Depends(current_active_user),
    user_manager: UserManager = Depends(get_user_manager)
) -> UserRead:
    try:
        user = await user_manager.update(new_password, user, request)
    except exceptions.InvalidPasswordException as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail={
                "code": ErrorCode.UPDATE_USER_INVALID_PASSWORD,
                "reason": e.reason,
            },
        )
    return UserRead.model_validate(user)