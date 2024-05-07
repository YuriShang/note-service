from typing import Optional, Annotated, Any, Callable
from datetime import datetime
import uuid

from pydantic import (
    BaseModel, 
    ConfigDict, 
    StringConstraints, 
    WrapValidator, 
    ValidationInfo
)
from pydantic_core import PydanticCustomError
from fastapi_users import schemas


def custom_error_msg(exc_factory: Callable[[str | None, Exception], Exception]) -> Any:
    def _validator(v: Any, next_: Any, ctx: ValidationInfo) -> Any:
        try:
            return next_(v, ctx)
        except Exception as e:
            raise exc_factory(ctx.field_name, e) from None
    return WrapValidator(_validator)


UsernameString = Annotated[
    str,
    StringConstraints(pattern=r"^[a-zA-Z]+$"),
    custom_error_msg(
        lambda field_name, _: PydanticCustomError(
            "Username field error",
            f"The field {field_name} must contain latin letters only",
        )
    ),
]


class UserRead(BaseModel):
    id: uuid.UUID
    username: str
    register_time: datetime
    password_set_time: datetime  

    model_config = ConfigDict(from_attributes=True)


class User(UserRead):
    hashed_password: str


class UserCreate(schemas.CreateUpdateDictModel):
    username: UsernameString
    password: str


class UserUpdate(schemas.CreateUpdateDictModel):
    username: Optional[UsernameString] = None
    password: Optional[str] = None
    
    
class UsernameUpdate(schemas.CreateUpdateDictModel):
    username: UsernameString


class PasswordUpdate(schemas.CreateUpdateDictModel):
    password: str