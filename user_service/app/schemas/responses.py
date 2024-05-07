from fastapi import status
from fastapi_users.router.common import ErrorCode, ErrorModel


username_change_responses = {
    status.HTTP_401_UNAUTHORIZED: {
        "description": "Missing token or inactive user.",
    },
    status.HTTP_400_BAD_REQUEST: {
        "model": ErrorModel,
        "content": {
            "application/json": {
                "examples": {
                    ErrorCode.UPDATE_USER_EMAIL_ALREADY_EXISTS: {
                        "summary": "A user with this email already exists.",
                        "value": {
                            "detail": ErrorCode.UPDATE_USER_EMAIL_ALREADY_EXISTS
                        },
                    },
                }
            }
        },
    },
}

password_change_responses = {
    status.HTTP_401_UNAUTHORIZED: {
        "description": "Missing token or inactive user.",
    },
    status.HTTP_400_BAD_REQUEST: {
        "model": ErrorModel,
        "content": {
            "application/json": {
                "examples": {
                    ErrorCode.UPDATE_USER_INVALID_PASSWORD: {
                        "summary": "Password validation failed.",
                        "value": {
                            "detail": {
                                "code": ErrorCode.UPDATE_USER_INVALID_PASSWORD,
                                "reason": "Password should be"
                                "at least 4 characters",
                            }
                        },
                    },
                }
            }
        },
    },
}